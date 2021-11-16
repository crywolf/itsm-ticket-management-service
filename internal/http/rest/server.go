package rest

import (
	"context"
	"net/http"
	"time"

	incidentsvc "github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident/service"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// Server is a http.Handler with dependencies
type Server struct {
	Addr      string
	URISchema string
	router    *httprouter.Router
	logger    *zap.SugaredLogger
	//authService             auth.Service
	//userService             usersvc.Service
	incidentService         incidentsvc.IncidentService
	inputPayloadConverters  jsonInputPayloadConverters
	presenters              jsonPresenters
	ExternalLocationAddress string
}

// Config contains server configuration and dependencies
type Config struct {
	Addr      string
	URISchema string
	Logger    *zap.SugaredLogger
	//AuthService             auth.Service
	//UserService             usersvc.Service
	IncidentService         incidentsvc.IncidentService
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
		Addr:      cfg.Addr,
		URISchema: URISchema,
		router:    r,
		logger:    cfg.Logger,
		//authService:             cfg.AuthService,
		//userService:             cfg.UserService,
		incidentService: cfg.IncidentService,
		//payloadValidator:        cfg.PayloadValidator,
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

	channelID := r.Header.Get("channel-id")
	ctx := context.WithValue(r.Context(), channelIDKey, channelID)

	//authToken := r.Header.Get("authorization")
	//ctx = context.WithValue(ctx, authKey, authToken)

	s.router.ServeHTTP(w, r.WithContext(ctx))
}

type channelIDType int

var channelIDKey channelIDType

// channelIDFromRequest returns channelID stored in request's context, if any.
func channelIDFromRequest(r *http.Request) (string, bool) {
	ch, ok := r.Context().Value(channelIDKey).(string)
	return ch, ok
}

// assertChannelID writes error message to response and returns error if channelID cannot be determined,
// otherwise it returns channelID
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
