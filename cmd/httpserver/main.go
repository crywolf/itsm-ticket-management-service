package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	fieldengineer "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer"
	fieldengineersvc "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer/service"
	incidentsvc "github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident/service"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
	externalusersvc "github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/external_user_service"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository/memory"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// for testing with in-memory repository
type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

func (c realClock) NowFormatted() types.DateTime {
	return types.DateTime(c.Now().Format(time.RFC3339))
}

func main() {
	l, _ := zap.NewProduction()
	defer func(l *zap.Logger) {
		_ = l.Sync()
	}(l)

	logger := l.Sugar()

	logger.Info("App starting...")

	loadEnvConfiguration()

	basicUserRepository := &memory.BasicUserRepositoryMemory{}
	// add test user - just for playing and testing
	basicUser := user.BasicUser{
		ExternalUserUUID: "83b231f2-5898-2658-70f4-5db03d1ccbc1",
		Name:             "Jan",
		Surname:          "Novák",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	basicUserUUID, err := basicUserRepository.AddBasicUser(context.Background(), "", basicUser)
	if err != nil {
		log.Fatal(err)
	}
	err = basicUser.SetUUID(basicUserUUID)
	if err != nil {
		log.Fatal(err)
	}

	clock := realClock{}
	fieldEngineerRepository := memory.NewFieldEngineerRepositoryMemory(clock, basicUserRepository)
	fieldEngineer := fieldengineer.FieldEngineer{
		BasicUser: basicUser,
	}
	err = fieldEngineer.CreatedUpdated.SetCreatedBy(basicUser)
	if err != nil {
		log.Fatal(err)
	}
	err = fieldEngineer.CreatedUpdated.SetUpdatedBy(basicUser)
	if err != nil {
		log.Fatal(err)
	}

	fieldEngID, err := fieldEngineerRepository.AddFieldEngineer(context.Background(), "", fieldEngineer)
	if err != nil {
		log.Fatal(err)
	}
	err = fieldEngineer.SetUUID(fieldEngID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n===> field eng UUID:", fieldEngID)

	fieldEngineerService := fieldengineersvc.NewFieldEngineerService(fieldEngineerRepository)

	incidentRepository := memory.NewIncidentRepositoryMemory(clock, basicUserRepository, fieldEngineerRepository)
	incidentService := incidentsvc.NewIncidentService(incidentRepository, fieldEngineerRepository)

	// External user service fetches user data from external service
	externalUserService, err := externalusersvc.NewService(basicUserRepository)
	if err != nil {
		logger.Fatalw("could not create external user service", "error", err)
	}

	// HTTP server
	server := rest.NewServer(rest.Config{
		Addr:                    viper.GetString("HTTPBindAddress"),
		URISchema:               "http://",
		Logger:                  logger,
		ExternalUserService:     externalUserService,
		IncidentService:         incidentService,
		FieldEngineerService:    fieldEngineerService,
		ExternalLocationAddress: viper.GetString("ExternalLocationAddress"),
	})

	srv := &http.Server{
		Addr:    server.Addr,
		Handler: server,
	}

	// Graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		// Trap sigterm or interrupt and gracefully shutdown the server
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		sig := <-sigint
		logger.Infof("Got signal: %s", sig)
		// We received a signal, shut down.

		// Gracefully shutdown the server, waiting max 'timeout' seconds for current operations to complete
		timeout := viper.GetInt("HTTPShutdownTimeoutInSeconds")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		logger.Info("Shutting down HTTP server...")
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Errorw("HTTP server Shutdown", "error", err)
		}
		logger.Info("HTTP server shutdown finished successfully")

		close(idleConnsClosed)
	}()

	// Start the server
	logger.Infof("Starting HTTP server at %s", server.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logger.Fatalw("HTTP server ListenAndServe", "error", err)
	}

	// Block until a signal is received and graceful shutdown completed.
	<-idleConnsClosed

	logger.Info("Exiting")
	_ = logger.Sync()

}
