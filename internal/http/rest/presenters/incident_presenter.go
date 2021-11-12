package presenters

import (
	"fmt"
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters/hypermedia"
	"go.uber.org/zap"
)

// NewIncidentPresenter creates an incident presentation service
func NewIncidentPresenter(logger *zap.SugaredLogger, serverAddr string) IncidentPresenter {
	return &incidentPresenter{
		BasePresenter: NewBasePresenter(logger, serverAddr),
	}
}

type incidentPresenter struct {
	*BasePresenter
}

func (p incidentPresenter) WriteIncident(w http.ResponseWriter, incident incident.Incident, hypermediaMapper hypermedia.Mapper) {
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

	incHypermedia := p.resourceToHypermediaLinks(hypermediaMapper, incident)
	incHypermedia["self"] = map[string]string{
		"href": hypermediaMapper.SelfLink(),
	}

	incResp := api.IncidentResponse{
		Incident: apiInc,
		Links:    incHypermedia,
	}

	p.encodeJSON(w, incResp)
}

func (p incidentPresenter) WriteIncidentList(w http.ResponseWriter, incidentList []incident.Incident, hypermediaMapper hypermedia.Mapper) {
	var apiList []api.IncidentResponse

	for _, inc := range incidentList {
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

		incHypermedia := p.resourceToHypermediaLinks(hypermediaMapper, inc)
		incHypermedia["self"] = map[string]string{
			"href": fmt.Sprintf("%s/%s/%s", p.serverAddr, "incidents", inc.UUID()),
		}

		incResp := api.IncidentResponse{
			Incident: apiInc,
			Links:    incHypermedia,
		}
		apiList = append(apiList, incResp)
	}

	listLinks := api.HypermediaLinks{
		"self": map[string]string{
			"href": hypermediaMapper.SelfLink(),
		},
	}

	resp := api.IncidentListResponse{
		Result: apiList,
		Links:  listLinks,
	}

	p.encodeJSON(w, resp)
}
