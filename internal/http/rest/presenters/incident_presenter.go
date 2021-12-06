package presenters

import (
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

func (p incidentPresenter) RenderIncident(w http.ResponseWriter, inc incident.Incident, hypermediaMapper hypermedia.Mapper) {
	incResp := api.IncidentResponse{
		Incident: p.convertIncidentToAPI(inc),
		Links:    p.resourceToHypermediaLinks(inc, hypermediaMapper, false),
	}

	p.renderJSON(w, incResp)
}

func (p incidentPresenter) RenderIncidentList(w http.ResponseWriter, incidentList repository.IncidentList, hypermediaMapper hypermedia.Mapper) {
	var apiList []api.IncidentResponse

	for _, inc := range incidentList.Result {
		incResp := api.IncidentResponse{
			Incident: p.convertIncidentToAPI(inc),
			Links:    p.resourceToHypermediaLinks(inc, hypermediaMapper, true),
		}
		apiList = append(apiList, incResp)
	}

	pagination := api.Pagination{
		Total: incidentList.Total,
		Size:  incidentList.Size,
		Page:  incidentList.Page,
	}

	resp := api.IncidentListResponse{
		Result:     apiList,
		Pagination: pagination,
		Links:      p.hypermediaListLinks(hypermediaMapper, incidentList.Pagination),
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
