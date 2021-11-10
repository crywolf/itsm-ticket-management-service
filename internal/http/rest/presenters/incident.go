package presenters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"go.uber.org/zap"
)

// HypermediaActions maps allowed incident actions to hypermedia actions
type HypermediaActions map[incident.AllowedAction]api.ActionLink

// Presenter provides REST responses
type Presenter interface {
	// WriteError replies to the request with the specified error message and HTTP code.
	// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
	// The error message should be plain text.
	WriteError(w http.ResponseWriter, error string, code int)

	WriteIncident(w http.ResponseWriter, incident incident.Incident, actions HypermediaActions, selfURI string)
	WriteIncidentList(w http.ResponseWriter, incidentList []incident.Incident, actions HypermediaActions, selfURI string)
}

// NewPresenter creates a presentation service
func NewPresenter(logger *zap.SugaredLogger, serverAddr string) Presenter {
	return &presenter{
		logger:     logger,
		serverAddr: serverAddr,
	}
}

type presenter struct {
	logger     *zap.SugaredLogger
	serverAddr string
}

func (p presenter) WriteError(w http.ResponseWriter, error string, code int) {
	p.sendErrorJSON(w, error, code)
}

func (p presenter) WriteIncident(w http.ResponseWriter, incident incident.Incident, actions HypermediaActions, selfURI string) {
	createdInfo := api.CreatedInfo{
		CreatedAt: incident.CreatedUpdated.CreatedAt().String(),
		CreatedBy: incident.CreatedUpdated.CreatedBy().String(),
	}

	updatedInfo := api.UpdatedInfo{
		UpdatedAt: incident.CreatedUpdated.UpdatedAt().String(),
		UpdatedBy: incident.CreatedUpdated.UpdatedBy().String(),
	}

	apiInc := api.Incident{
		UUID:             incident.UUID().String(),
		ExternalID:       incident.ExternalID,
		ShortDescription: incident.ShortDescription,
		Description:      incident.Description,
		State:            incident.State(),
		CreatedUpdated: api.CreatedUpdated{
			CreatedInfo: createdInfo,
			UpdatedInfo: updatedInfo,
		},
	}

	incHypermedia := api.HypermediaLinks{}

	allowedActions := incident.AllowedActions()
	for _, action := range allowedActions {
		link := actions[action]
		href := strings.ReplaceAll(link.Href, "{uuid}", incident.UUID().String())

		incHypermedia[link.Name] = map[string]string{
			"href": href,
		}
	}
	incHypermedia["self"] = map[string]string{
		"href": selfURI,
	}

	incResp := api.IncidentResponse{
		Incident: apiInc,
		Links:    incHypermedia,
	}

	p.encodeJSON(w, incResp)
}

func (p presenter) WriteIncidentList(w http.ResponseWriter, incList []incident.Incident, actions HypermediaActions, selfURI string) {
	var apiList []api.IncidentResponse

	for _, inc := range incList {
		createdInfo := api.CreatedInfo{
			CreatedAt: inc.CreatedUpdated.CreatedAt().String(),
			CreatedBy: inc.CreatedUpdated.CreatedBy().String(),
		}

		updatedInfo := api.UpdatedInfo{
			UpdatedAt: inc.CreatedUpdated.UpdatedAt().String(),
			UpdatedBy: inc.CreatedUpdated.UpdatedBy().String(),
		}

		apiInc := api.Incident{
			UUID:             inc.UUID().String(),
			ExternalID:       inc.ExternalID,
			ShortDescription: inc.ShortDescription,
			Description:      inc.Description,
			State:            inc.State(),
			CreatedUpdated: api.CreatedUpdated{
				CreatedInfo: createdInfo,
				UpdatedInfo: updatedInfo,
			},
		}

		incHypermedia := api.HypermediaLinks{}

		allowedActions := inc.AllowedActions()
		for _, action := range allowedActions {
			link := actions[action]
			href := strings.ReplaceAll(link.Href, "{uuid}", inc.UUID().String())

			incHypermedia[link.Name] = map[string]string{
				"href": href,
			}
		}
		incHypermedia["self"] = map[string]string{
			"href": fmt.Sprintf("%s/%s", selfURI, inc.UUID()),
		}

		incResp := api.IncidentResponse{
			Incident: apiInc,
			Links:    incHypermedia,
		}
		apiList = append(apiList, incResp)
	}

	links := api.HypermediaLinks{
		"self": map[string]string{
			"href": selfURI,
		},
	}

	resp := api.IncidentListResponse{
		Result: apiList,
		Links:  links,
	}

	p.encodeJSON(w, resp)
}

// sendErrorJSON replies to the request with the specified error message and HTTP code.
// It encodes error string as JSON object {"error":"error_string"} and sets correct header.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
// The error message should be plain text.
func (p presenter) sendErrorJSON(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errorJSON, _ := json.Marshal(error)
	_, _ = fmt.Fprintf(w, `{"error":%s}`+"\n", errorJSON)
}

// encodeJSON encodes 'v' to JSON and writes it to the 'w'. Also sets correct Content-Type header.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
func (p presenter) encodeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		eMsg := "could not encode JSON response"
		p.logger.Errorw(eMsg, "error", err)
		p.WriteError(w, eMsg, http.StatusInternalServerError)
		return
	}
}
