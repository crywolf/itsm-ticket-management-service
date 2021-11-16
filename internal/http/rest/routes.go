package rest

import (
	"embed"
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
	"github.com/go-openapi/runtime/middleware"
)

//go:embed api/swagger.yaml
var swaggerFS embed.FS

func (s *Server) registerRoutes() {
	s.registerIncidentRoutes()

	// API documentation
	opts := middleware.RedocOpts{Path: "/docs", SpecURL: "/swagger.yaml", Title: "Ticket management service API documentation"}
	docsHandler := middleware.Redoc(opts, nil)
	// handlers for API documentation
	s.router.Handler(http.MethodGet, "/docs", docsHandler)
	s.router.Handler(http.MethodGet, "/swagger.yaml", http.FileServer(http.FS(swaggerFS)))

	// default Not Found handler
	s.router.NotFound = http.HandlerFunc(s.JSONNotFoundError)
}

// JSONNotFoundError replies to the request with the 404 page not found general error message
// in JSON format and sets correct header and HTTP code
func (s Server) JSONNotFoundError(w http.ResponseWriter, _ *http.Request) {
	s.presenters.base.RenderError(w, "", presenters.NewErrorf(http.StatusNotFound, "404 page not found"))
}
