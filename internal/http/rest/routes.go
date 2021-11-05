package rest

import (
	"embed"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

//go:embed api/swagger.yaml
var swaggerFS embed.FS

func (s *Server) routes() {
	router := s.router

	// incidents
	router.GET("/incidents/:id", s.GetIncident())
	//	router.GET("/incidents", s.ListIncidents)

	router.POST("/incidents", s.CreateIncident())

	// API documentation
	opts := middleware.RedocOpts{Path: "/docs", SpecURL: "/swagger.yaml", Title: "Ticket management service API documentation"}
	docsHandler := middleware.Redoc(opts, nil)
	// handlers for API documentation
	router.Handler(http.MethodGet, "/docs", docsHandler)
	router.Handler(http.MethodGet, "/swagger.yaml", http.FileServer(http.FS(swaggerFS)))

	// default Not Found handler
	router.NotFound = http.HandlerFunc(s.JSONNotFoundError)
}

// JSONNotFoundError replies to the request with the 404 page not found general error message
// in JSON format and sets correct header and HTTP code
func (s Server) JSONNotFoundError(w http.ResponseWriter, _ *http.Request) {
	s.presenter.WriteError(w, "404 page not found", http.StatusNotFound)
}
