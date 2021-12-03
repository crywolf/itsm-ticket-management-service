package presenters

import (
	"fmt"
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters/hypermedia"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
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

func (p incidentPresenter) RenderIncident(w http.ResponseWriter, incident incident.Incident, hypermediaMapper hypermedia.Mapper) {
	apiInc := p.convertIncidentToAPI(incident)

	incHypermedia := p.resourceToHypermediaLinks(hypermediaMapper, incident)
	incHypermedia["self"] = map[string]string{
		"href": hypermediaMapper.SelfLink(),
	}

	incResp := api.IncidentResponse{
		Incident: apiInc,
		Links:    incHypermedia,
	}

	p.renderJSON(w, incResp)
}

func (p incidentPresenter) RenderIncidentList(w http.ResponseWriter, incidentList repository.IncidentList, listRoute string, hypermediaMapper hypermedia.Mapper) {
	var apiList []api.IncidentResponse

	for _, inc := range incidentList.Result {
		apiInc := p.convertIncidentToAPI(inc)

		incHypermedia := p.resourceToHypermediaLinks(hypermediaMapper, inc)
		incHypermedia["self"] = map[string]string{
			"href": fmt.Sprintf("%s%s/%s", p.serverAddr, listRoute, inc.UUID()),
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
		Total:  incidentList.Total,
		Size:   incidentList.Size,
		Page:   incidentList.Page,
		Links:  listLinks,
	}

	p.renderJSON(w, resp)
}

func (p incidentPresenter) convertIncidentToAPI(inc incident.Incident) api.Incident {
	createdInfo := api.CreatedInfo{
		CreatedAt: inc.CreatedUpdated.CreatedAt().String(),
		CreatedBy: inc.CreatedUpdated.CreatedBy().String(),
	}

	updatedInfo := api.UpdatedInfo{
		UpdatedAt: inc.CreatedUpdated.UpdatedAt().String(),
		UpdatedBy: inc.CreatedUpdated.UpdatedBy().String(),
	}

	var timelogUUIDs []api.UUID
	for _, timelog := range inc.Timelogs {
		timelogUUIDs = append(timelogUUIDs, api.UUID(timelog))
	}

	apiInc := api.Incident{
		UUID:             inc.UUID().String(),
		Number:           inc.Number,
		ExternalID:       inc.ExternalID,
		ShortDescription: inc.ShortDescription,
		Description:      inc.Description,
		State:            inc.State(),
		Timelogs:         timelogUUIDs,
		CreatedUpdated: api.CreatedUpdated{
			CreatedInfo: createdInfo,
			UpdatedInfo: updatedInfo,
		},
	}
	return apiInc
}
