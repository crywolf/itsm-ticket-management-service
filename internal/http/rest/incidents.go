package rest

import (
	"net/http"
	"net/url"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters/hypermedia"
	"github.com/julienschmidt/httprouter"
)

func (s Server) registerIncidentRoutes() {
	s.router.POST("/incidents", s.CreateIncident())
	s.router.GET("/incidents/:id", s.GetIncident())
	s.router.GET("/incidents", s.ListIncidents())
}

// swagger:route POST /incidents incidents CreateIncident
// Creates a new incident
// responses:
//	201: incidentCreatedResponse
//	400: errorResponse400
//	401: errorResponse401
//	403: errorResponse403
//	409: errorResponse409

// CreateIncident returns handler for creating single incident
func (s *Server) CreateIncident() func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		incPayload, err := s.inputPayloadConverters.incident.IncidentParamsFromBody(r)
		if err != nil {
			s.logger.Warnw("CreateIncident handler failed", "error", err)
			s.presenters.incident.RenderError(w, "", err)
			return
		}

		channelID, err := s.assertChannelID(w, r)
		if err != nil {
			return
		}

		actorUser, err := s.actorFromRequest(r)
		if err != nil {
			s.logger.Errorw("CreateIncident handler failed", "error", err)
			s.presenters.incident.RenderError(w, "", err)
			return
		}

		newID, err := s.incidentService.CreateIncident(r.Context(), channelID, actorUser, incPayload)
		if err != nil {
			s.logger.Errorw("CreateIncident handler failed", "error", err)
			s.presenters.incident.RenderError(w, "", err)
			return
		}

		s.presenters.incident.RenderLocationHeader(w, listIncidentsRoute, newID)
	}
}

// swagger:route GET /incidents/{uuid} incidents GetIncident
// Returns a single incident from the repository
// responses:
//	200: incidentResponse
//	400: errorResponse400
//  401: errorResponse401
//  403: errorResponse403
//	404: errorResponse404

// GetIncident returns handler for getting single incident
func (s *Server) GetIncident() func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		if id == "" {
			err := presenters.NewErrorf(http.StatusBadRequest, "malformed URL: missing resource ID param")
			s.logger.Errorw("GetIncident handler failed", "error", err)
			s.presenters.base.RenderError(w, "", err)
			return
		}

		channelID, err := s.assertChannelID(w, r)
		if err != nil {
			return
		}

		actorUser, err := s.actorFromRequest(r)
		if err != nil {
			s.logger.Errorw("GetIncident handler failed", "error", err)
			s.presenters.base.RenderError(w, "", err)
			return
		}

		inc, err := s.incidentService.GetIncident(r.Context(), channelID, actorUser, ref.UUID(id))
		if err != nil {
			s.logger.Errorw("GetIncident handler failed", "ID", id, "error", err)
			s.presenters.base.RenderError(w, "incident not found", err)
			return
		}

		hypermediaMapper := NewIncidentHypermediaMapper(s.ExternalLocationAddress, r.URL, actorUser)
		s.presenters.incident.RenderIncident(w, inc, hypermediaMapper)
	}
}

// swagger:route GET /incidents incidents ListIncidents
// Returns a list of incidents
// responses:
//	200: incidentListResponse
//	400: errorResponse400
//  401: errorResponse401
//  403: errorResponse403
const listIncidentsRoute = "/incidents"

// ListIncidents returns handler for listing incidents
func (s *Server) ListIncidents() func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		channelID, err := s.assertChannelID(w, r)
		if err != nil {
			return
		}

		actorUser, err := s.actorFromRequest(r)
		if err != nil {
			s.logger.Errorw("ListIncidents handler failed", "error", err)
			s.presenters.base.RenderError(w, "", err)
			return
		}

		paginationParams, err := s.PaginationParams(r, actorUser)
		if err != nil {
			s.presenters.base.RenderError(w, "", err)
			return
		}

		list, err := s.incidentService.ListIncidents(r.Context(), channelID, actorUser, paginationParams)
		if err != nil {
			s.logger.Errorw("ListIncidents handler failed", "error", err)
			s.presenters.base.RenderError(w, "", err)
			return
		}

		hypermediaMapper := NewIncidentHypermediaMapper(s.ExternalLocationAddress, r.URL, actorUser)
		s.presenters.incident.RenderIncidentList(w, list, hypermediaMapper)
	}
}

// IncidentHypermediaMapper implements hypermedia mapping functionality for incident resource
type IncidentHypermediaMapper struct {
	*hypermedia.BaseHypermediaMapper
}

// NewIncidentHypermediaMapper returns new hypermedia mapper for incident resource
func NewIncidentHypermediaMapper(serverAddr string, currentURL *url.URL, actor actor.Actor) IncidentHypermediaMapper {
	return IncidentHypermediaMapper{
		BaseHypermediaMapper: hypermedia.NewBaseHypermedia(serverAddr, currentURL, actor),
	}
}

// RoutesToHypermediaActionLinks maps domain object actions to hypermedia action links
func (h IncidentHypermediaMapper) RoutesToHypermediaActionLinks() hypermedia.ActionLinks {
	acts := hypermedia.ActionLinks{}

	// TODO implement as a method on BaseHypermediaMapper
	acts[incident.ActionCancel.String()] = api.ActionLink{Name: "CancelIncident", Href: h.ServerAddr() + cancelIncidentRoute}
	acts[incident.ActionStartWorking.String()] = api.ActionLink{Name: "IncidentStartWorking", Href: h.ServerAddr() + incidentStartWorkingRoute}

	return acts
}

// TODO implement routes - they are just for testing at the moment
const cancelIncidentRoute = "/incidents/{uuid}/cancel"

const incidentStartWorkingRoute = "/incidents/{uuid}/start_working"
