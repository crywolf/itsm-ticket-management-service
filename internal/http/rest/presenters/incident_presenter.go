package presenters

import (
	"fmt"
	"net/http"
	"strconv"

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

	listLinks := api.HypermediaLinks{}
	listLinks.AppendSelfLink(hypermediaMapper.SelfLink())

	// TODO move to hypermedia object
	query := hypermediaMapper.RequestURL().Query()
	url := *hypermediaMapper.RequestURL()
	first := url
	last := url
	prev := url
	next := url

	if incidentList.First == 1 {
		query.Del("page")
	} else {
		query.Set("page", strconv.Itoa(incidentList.First))
	}
	first.RawQuery = query.Encode()

	if incidentList.Last == 1 {
		query.Del("page")
	} else {
		query.Set("page", strconv.Itoa(incidentList.Last))
	}
	last.RawQuery = query.Encode()

	prevString := ""
	if incidentList.Prev == 1 {
		query.Del("page")
		prev.RawQuery = query.Encode()
		prevString = fmt.Sprintf("%s%s", p.serverAddr, prev.String())
	} else if incidentList.Prev > 1 {
		query.Set("page", strconv.Itoa(incidentList.Prev))
		prev.RawQuery = query.Encode()
		prevString = fmt.Sprintf("%s%s", p.serverAddr, prev.String())
	}

	nextString := ""
	if incidentList.Next > 0 {
		query.Set("page", strconv.Itoa(incidentList.Next))
		next.RawQuery = query.Encode()
		nextString = fmt.Sprintf("%s%s", p.serverAddr, next.String())
	}

	pagination := api.Pagination{
		Total: incidentList.Total,
		Size:  incidentList.Size,
		Page:  incidentList.Page,
		// TODO move it to _links{} in hypermedia | in IncidentListResponse Links: listLinks
		First: fmt.Sprintf("%s%s", p.serverAddr, first.String()),
		Last:  fmt.Sprintf("%s%s", p.serverAddr, last.String()),
		Prev:  prevString,
		Next:  nextString,
	}

	resp := api.IncidentListResponse{
		Result:     apiList,
		Pagination: pagination,
		Links:      listLinks,
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
