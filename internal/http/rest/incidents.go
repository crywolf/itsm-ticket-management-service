package rest

import (
	"context"
	"net/http"
	"net/url"

	fieldengineersvc "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer/service"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters/hypermedia"
	"github.com/julienschmidt/httprouter"
)

func (s Server) registerIncidentRoutes() {
	s.router.POST("/incidents", s.CreateIncident())
	s.router.PATCH("/incidents/:id", s.UpdateIncident())
	s.router.GET("/incidents/:id", s.GetIncident())
	s.router.GET("/incidents", s.ListIncidents())
	s.router.POST("/incidents/:id/start_working", s.IncidentStartWorking())
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
		incPayload, err := s.inputPayloadConverters.incident.IncidentCreateParamsFromBody(r)
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

		s.presenters.incident.RenderCreatedHeader(w, listIncidentsRoute, newID)
	}
}

// swagger:route PATCH /incidents/{uuid} incidents UpdateIncident
// Updates specified incident
// responses:
//	204: incidentNoContentResponse
//	400: errorResponse400
//	401: errorResponse401
//  403: errorResponse403
//	404: errorResponse404

// UpdateIncident returns handler for creating single incident
func (s *Server) UpdateIncident() func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		if id == "" {
			err := presenters.NewErrorf(http.StatusBadRequest, "malformed URL: missing resource ID param")
			s.logger.Errorw("UpdateIncident handler failed", "error", err)
			s.presenters.base.RenderError(w, "", err)
			return
		}

		incPayload, err := s.inputPayloadConverters.incident.IncidentUpdateParamsFromBody(r)
		if err != nil {
			s.logger.Warnw("UpdateIncident handler failed", "error", err)
			s.presenters.incident.RenderError(w, "", err)
			return
		}

		channelID, err := s.assertChannelID(w, r)
		if err != nil {
			return
		}

		actorUser, err := s.actorFromRequest(r)
		if err != nil {
			s.logger.Errorw("UpdateIncident handler failed", "error", err)
			s.presenters.incident.RenderError(w, "", err)
			return
		}

		newID, err := s.incidentService.UpdateIncident(r.Context(), channelID, actorUser, ref.UUID(id), incPayload)
		if err != nil {
			s.logger.Errorw("UpdateIncident handler failed", "error", err)
			s.presenters.incident.RenderError(w, "", err)
			return
		}
		s.presenters.incident.RenderNoContentHeader(w, listIncidentsRoute, newID)
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
const getIncidentRoute = "/incidents/{uuid}"

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

		hypermediaMapper := NewIncidentHypermediaMapper(r.Context(), channelID, s.ExternalLocationAddress, r.URL, actorUser, s.fieldEngineerService)
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

		hypermediaMapper := NewIncidentHypermediaMapper(r.Context(), channelID, s.ExternalLocationAddress, r.URL, actorUser, s.fieldEngineerService)
		s.presenters.incident.RenderIncidentList(w, list, hypermediaMapper)
	}
}

// swagger:route POST /incidents/{uuid}/start_working incidents IncidentStartWorking
// Starts working on incident by field engineer
// responses:
//	204: incidentNoContentResponse
//	400: errorResponse400
//	401: errorResponse401
//  403: errorResponse403
const incidentStartWorkingRoute = "/incidents/{uuid}/start_working"

// IncidentStartWorking returns handler for start working action
func (s *Server) IncidentStartWorking() func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		incID := params.ByName("id")
		if incID == "" {
			err := presenters.NewErrorf(http.StatusBadRequest, "malformed URL: missing resource ID param")
			s.logger.Errorw("IncidentStartWorking handler failed", "error", err)
			s.presenters.base.RenderError(w, "", err)
			return
		}

		payload, err := s.inputPayloadConverters.incident.IncidentStartWorkingParamsFromBody(r)
		if err != nil {
			s.logger.Warnw("IncidentStartWorking handler failed", "error", err)
			s.presenters.incident.RenderError(w, "", err)
			return
		}

		channelID, err := s.assertChannelID(w, r)
		if err != nil {
			return
		}

		actorUser, err := s.actorFromRequest(r)
		if err != nil {
			s.logger.Errorw("IncidentStartWorking handler failed", "error", err)
			s.presenters.incident.RenderError(w, "", err)
			return
		}

		err = s.incidentService.StartWorking(r.Context(), channelID, actorUser, ref.UUID(incID), payload)
		if err != nil {
			s.logger.Errorw("IncidentStartWorking handler failed", "error", err)
			s.presenters.incident.RenderError(w, "", err)
			return
		}

		s.presenters.incident.RenderNoContentHeader(w, listIncidentsRoute, ref.UUID(incID))
	}
}

// IncidentHypermediaMapper implements hypermedia mapping functionality for incident resource
type IncidentHypermediaMapper struct {
	ctx       context.Context
	channelID ref.ChannelID
	feSvc     fieldengineersvc.FieldEngineerService
	*hypermedia.BaseHypermediaMapper
}

// NewIncidentHypermediaMapper returns new hypermedia mapper for incident resource
func NewIncidentHypermediaMapper(ctx context.Context, channelID ref.ChannelID, serverAddr string, currentURL *url.URL, actor actor.Actor, feSvc fieldengineersvc.FieldEngineerService) IncidentHypermediaMapper {
	return IncidentHypermediaMapper{
		ctx:                  ctx,
		channelID:            channelID,
		feSvc:                feSvc,
		BaseHypermediaMapper: hypermedia.NewBaseHypermedia(serverAddr, currentURL, actor),
	}
}

// Ctx ...
func (h IncidentHypermediaMapper) Ctx() context.Context {
	return h.ctx
}

// ChannelID ...
func (h IncidentHypermediaMapper) ChannelID() ref.ChannelID {
	return h.channelID
}

// FieldEngineerSvc ...
func (h IncidentHypermediaMapper) FieldEngineerSvc() fieldengineersvc.FieldEngineerService {
	return h.feSvc
}

// RoutesToHypermediaActionLinks maps domain object actions to hypermedia action links
func (h IncidentHypermediaMapper) RoutesToHypermediaActionLinks() hypermedia.ActionLinks {
	links := hypermedia.NewActionLinks(h.BaseHypermediaMapper)

	links.Add(incident.ActionCancel.String(), "CancelIncident", cancelIncidentRoute)
	links.Add(incident.ActionStartWorking.String(), "IncidentStartWorking", incidentStartWorkingRoute)

	return links
}

// TODO this route will be in used in hypermedia
const getBasicUserRoute = "/basic_user/{uuid}"

// TODO implement routes - they are just for testing at the moment
const cancelIncidentRoute = "/incidents/{uuid}/cancel"
