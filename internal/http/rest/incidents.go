package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

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
		defer func() { _ = r.Body.Close() }()
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.logger.Errorw("could not read request body", "error", err)
			s.presenter.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var incPayload api.CreateIncidentParams

		err = json.Unmarshal(reqBody, &incPayload)
		if err != nil {
			eMsg := "could not decode JSON from request"
			s.logger.Warnw("CreateIncident handler failed", "reason", eMsg, "error", err)
			s.presenter.WriteError(w, fmt.Sprintf("%s: %s", eMsg, err.Error()), http.StatusBadRequest)
			return
		}

		channelID, err := s.assertChannelID(w, r)
		if err != nil {
			return
		}

		//user, ok := s.UserInfoFromRequest(r)
		//if !ok {
		//	eMsg := "could not get invoking user from context"
		//	s.logger.Error(eMsg)
		//	s.presenter.WriteError(w, eMsg, http.StatusInternalServerError)
		//	return
		//}

		newID, err := s.incidentService.CreateIncident(r.Context(), channelID, incPayload)
		if err != nil {
			var httpError *repository.Error
			if errors.As(err, &httpError) {
				s.logger.Warnw("CreateIncident handler failed", "error", err)
				s.presenter.WriteError(w, err.Error(), httpError.StatusCode())
				return
			}

			s.logger.Errorw("CreateIncident handler failed", "error", err)
			s.presenter.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resourceURI := fmt.Sprintf("%s/%s/%s", s.ExternalLocationAddress, "incidents", newID)

		w.Header().Set("Location", resourceURI)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
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
			eMsg := "malformed URL: missing resource ID param"
			s.logger.Warnw("GetIncident handler failed", "error", eMsg)
			s.presenter.WriteError(w, eMsg, http.StatusBadRequest)
			return
		}

		channelID, err := s.assertChannelID(w, r)
		if err != nil {
			return
		}

		inc, err := s.incidentService.GetIncident(r.Context(), channelID, ref.UUID(id))
		if err != nil {
			var httpError *repository.Error
			if errors.As(err, &httpError) {
				s.logger.Warnw("GetIncident handler failed", "ID", id, "error", err)
				s.presenter.WriteError(w, err.Error(), httpError.StatusCode())
				return
			}

			s.logger.Errorw("GetIncident handler failed", "ID", id, "error", err)
			s.presenter.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		self := s.ExternalLocationAddress + r.URL.String()

		s.presenter.WriteIncident(w, inc, incidentRoutesToHypermedia(s.ExternalLocationAddress), self)
	}
}

// swagger:route GET /incidents incidents ListIncidents
// Returns a list of incidents
// responses:
//	200: incidentListResponse
//	400: errorResponse400
//  401: errorResponse401
//  403: errorResponse403

// ListIncidents returns handler for listing incidents
func (s *Server) ListIncidents() func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		channelID, err := s.assertChannelID(w, r)
		if err != nil {
			return
		}

		list, err := s.incidentService.ListIncidents(r.Context(), channelID)
		if err != nil {
			var httpError *repository.Error
			if errors.As(err, &httpError) {
				s.logger.Error("Repository error", zap.Error(err))
				s.presenter.WriteError(w, err.Error(), httpError.StatusCode())
				return
			}

			s.logger.Error("ListIncidents handler failed", zap.Error(err))
			s.presenter.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		self := s.ExternalLocationAddress + r.URL.String()

		s.presenter.WriteIncidentList(w, list, incidentRoutesToHypermedia(s.ExternalLocationAddress), self)
	}
}

func incidentRoutesToHypermedia(serverAddr string) presenters.HypermediaActions {
	acts := presenters.HypermediaActions{}

	acts[incident.ActionCancel] = api.ActionLink{Name: "CancelIncident", Href: serverAddr + cancelIncidentRoute}
	acts[incident.ActionStartWorking] = api.ActionLink{Name: "IncidentStartWorking", Href: serverAddr + incidentStartWorkingRoute}

	return acts
}

const cancelIncidentRoute = "/incidents/{uuid}/cancel"

const incidentStartWorkingRoute = "/incidents/{uuid}/start_working"
