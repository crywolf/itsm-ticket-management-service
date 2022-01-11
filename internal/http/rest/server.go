package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	fieldengineersvc "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer/service"
	incidentsvc "github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident/service"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	externalusersvc "github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/external_user_service"
	converters "github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/input_converters"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// Server is a http.Handler with dependencies
type Server struct {
	Addr                    string
	URISchema               string
	router                  *httprouter.Router
	logger                  *zap.SugaredLogger
	clock                   domain.Clock
	externalUserService     externalusersvc.Service
	incidentService         incidentsvc.IncidentService
	fieldEngineerService    fieldengineersvc.FieldEngineerService
	inputPayloadConverters  jsonInputPayloadConverters
	presenters              jsonPresenters
	ExternalLocationAddress string
}

// Config contains server configuration and dependencies
type Config struct {
	Addr                    string
	URISchema               string
	Logger                  *zap.SugaredLogger
	Clock                   domain.Clock
	ExternalUserService     externalusersvc.Service
	IncidentService         incidentsvc.IncidentService
	FieldEngineerService    fieldengineersvc.FieldEngineerService
	ExternalLocationAddress string
}

// NewServer creates new server with the necessary dependencies
func NewServer(cfg Config) *Server {
	r := httprouter.New()

	URISchema := "http://"
	if cfg.URISchema != "" {
		URISchema = cfg.URISchema
	}

	s := &Server{
		Addr:                    cfg.Addr,
		URISchema:               URISchema,
		router:                  r,
		logger:                  cfg.Logger,
		clock:                   cfg.Clock,
		externalUserService:     cfg.ExternalUserService,
		incidentService:         cfg.IncidentService,
		fieldEngineerService:    cfg.FieldEngineerService,
		ExternalLocationAddress: cfg.ExternalLocationAddress,
	}
	s.registerInputConverters()
	s.registerPresenters()
	s.registerRoutes()

	return s
}

// ServeHTTP makes the server implement the http.Handler interface
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Infow(r.Method,
		"time", time.Now().Format(time.RFC3339),
		"url", r.URL.String(),
	)

	// add channelID and authToken to request's context
	sChannelID := r.Header.Get("channel-id")
	ctx := context.WithValue(r.Context(), channelIDKey, sChannelID)

	authToken := r.Header.Get("authorization")
	ctx = context.WithValue(ctx, authKey, authToken)

	r = r.WithContext(ctx)

	// make sure the authToken was sent in the header, otherwise end here and render error to the client
	if _, err := s.assertAuthToken(w, r); err != nil {
		return
	}

	// make sure the channelID was sent in the header, otherwise end here and render error to the client
	channelID, err := s.assertChannelID(w, r)
	if err != nil {
		return
	}

	// get Actor from the external user service and add it to the request
	actorUser, err := s.externalUserService.ActorFromRequest(ctx, authToken, channelID, r.Header.Get("on_behalf"))
	if err != nil {
		s.logger.Errorw("externalUserService.ActorFromRequest failed:", "error", err)
		s.presenters.base.RenderError(w, "", err)
		return
	}
	ctx = context.WithValue(ctx, userKey, &actorUser)

	s.router.ServeHTTP(w, r.WithContext(ctx))
}

type userKeyType int

var userKey userKeyType

// ActorFromRequest returns the Actor user stored in request's context if any.
func (s Server) actorFromRequest(r *http.Request) (actor.Actor, error) {
	act, ok := r.Context().Value(userKey).(*actor.Actor)
	if !ok {
		err := presenters.NewErrorf(http.StatusInternalServerError, "could not get actor from context")
		return actor.Actor{}, err
	}

	return *act, nil
}

type channelIDType int

var channelIDKey channelIDType

// channelIDFromRequest returns channelID stored in request's context, if any.
func channelIDFromRequest(r *http.Request) (string, bool) {
	ch, ok := r.Context().Value(channelIDKey).(string)
	return ch, ok
}

// assertChannelID writes error message to response and returns error if channelID cannot be determined,
// otherwise it returns channelID.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
func (s Server) assertChannelID(w http.ResponseWriter, r *http.Request) (ref.ChannelID, error) {
	channelID, ok := channelIDFromRequest(r)
	if !ok {
		err := presenters.NewErrorf(http.StatusInternalServerError, "could not get channel ID from context")
		s.logger.Errorw("assertChannelID", "error", err)
		s.presenters.base.RenderError(w, "cannot determine channel ID", err)
		return "", err
	}

	if channelID == "" {
		err := presenters.NewErrorf(http.StatusUnauthorized, "empty channel ID in context")
		s.logger.Errorw("assertChannelID", "error", err)
		s.presenters.base.RenderError(w, "'channel-id' header missing or invalid", err)
		return "", err
	}

	return ref.ChannelID(channelID), nil
}

type authType int

var authKey authType

// assertAuthToken writes error message to response and returns error if authorization token cannot be determined,
// otherwise it returns authorization token.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
func (s Server) assertAuthToken(w http.ResponseWriter, r *http.Request) (string, error) {
	authToken, ok := authTokenFromRequest(r)
	if !ok {
		err := presenters.NewErrorf(http.StatusInternalServerError, "could not get authorization token from context")
		s.logger.Errorw("assertAuthToken", "error", err)
		s.presenters.base.RenderError(w, "cannot determine authorization token", err)
		return "", err
	}

	if authToken == "" {
		err := presenters.NewErrorf(http.StatusUnauthorized, "empty authorization token in context")
		s.logger.Errorw("assertAuthToken", "error", err)
		s.presenters.base.RenderError(w, "'authorization' header missing or invalid", err)
		return "", err
	}

	return authToken, nil
}

// authTokenFromRequest returns authorization token stored in request's context, if any.
func authTokenFromRequest(r *http.Request) (string, bool) {
	ch, ok := r.Context().Value(authKey).(string)
	return ch, ok
}

// PaginationParams parses request query and returns params with information about requested page and items per page to be displayed
func (s Server) PaginationParams(r *http.Request, actorUser actor.Actor) (converters.PaginationParams, error) {
	return converters.NewPaginationParams(r, actorUser)
}
