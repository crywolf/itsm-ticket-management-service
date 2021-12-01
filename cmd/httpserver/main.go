package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	incidentsvc "github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident/service"
	usersvc "github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/service"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository/memory"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

func main() {
	l, _ := zap.NewProduction()
	defer func(l *zap.Logger) {
		_ = l.Sync()
	}(l)

	logger := l.Sugar()

	logger.Info("App starting...")

	loadEnvConfiguration()

	r := &memory.IncidentRepositoryMemory{
		Clock: realClock{},
	}
	incidentService := incidentsvc.NewIncidentService(r)

	// User service fetches user data from external service
	userService, err := usersvc.NewService()
	if err != nil {
		logger.Fatal("could not create user service", zap.Error(err))
	}

	// HTTP server
	server := rest.NewServer(rest.Config{
		Addr:                    viper.GetString("HTTPBindAddress"),
		URISchema:               "http://",
		Logger:                  logger,
		UserService:             userService,
		IncidentService:         incidentService,
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
