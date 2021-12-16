package presenters

import (
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/embedded"
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

func (p incidentPresenter) RenderIncident(w http.ResponseWriter, inc incident.Incident, hypermediaMapper hypermedia.IncidentMapper) {
	// TODO improve this embedded mapper programming interface - make it an object
	var embeddedMappings []hypermedia.EmbeddedResourceMapping

	if inc.FieldEngineerID != nil {
		feSvc := hypermediaMapper.FieldEngineerSvc()
		fe, err := feSvc.GetFieldEngineer(hypermediaMapper.Ctx(), hypermediaMapper.ChannelID(), hypermediaMapper.Actor(), *inc.FieldEngineerID)
		if err != nil {
			err = WrapErrorf(err, http.StatusInternalServerError, "error rendering embedded resource")
			p.RenderError(w, "", err)
		}
		embeddedFieldEngineer := api.NewEmbeddedFieldEngineer(fe)
		mappingFE := *hypermedia.EmbeddedResourcesMappingDefinition[embedded.FieldEngineer].AddResource(embeddedFieldEngineer)
		embeddedMappings = append(embeddedMappings, mappingFE)
	}

	embeddedCreatedBy := api.NewEmbeddedBasicUser(inc.CreatedUpdated.CreatedBy())
	mappingCreatedBy := *hypermedia.EmbeddedResourcesMappingDefinition[embedded.CreatedBy].AddResource(embeddedCreatedBy)
	//mappingCreatedBy := *hypermedia.EmbeddedResourcesMappingDefinition[embedded.UpdatedBy].AddResource(embeddedCreatedBy)

	embeddedMappings = append(embeddedMappings, mappingCreatedBy)

	incResp := api.IncidentResponse{
		Incident: p.convertIncidentToAPI(inc),
		Links:    p.resourceToHypermediaLinks(inc, hypermediaMapper, false),
		Embedded: p.resourceToEmbeddedField(inc, embeddedMappings, hypermediaMapper),
	}

	p.renderJSON(w, incResp)
}

func (p incidentPresenter) RenderIncidentList(w http.ResponseWriter, incidentList repository.IncidentList, hypermediaMapper hypermedia.IncidentMapper) {
	var apiList []api.IncidentResponse

	for _, inc := range incidentList.Result {
		var embeddedMappings []hypermedia.EmbeddedResourceMapping

		if inc.FieldEngineerID != nil {
			feSvc := hypermediaMapper.FieldEngineerSvc()
			fe, err := feSvc.GetFieldEngineer(hypermediaMapper.Ctx(), hypermediaMapper.ChannelID(), hypermediaMapper.Actor(), *inc.FieldEngineerID)
			if err != nil {
				err = WrapErrorf(err, http.StatusInternalServerError, "error rendering embedded resource")
				p.RenderError(w, "", err)
				return
			}
			embeddedFieldEngineer := api.NewEmbeddedFieldEngineer(fe)
			mappingFE := *hypermedia.EmbeddedResourcesMappingDefinition[embedded.FieldEngineer].AddResource(embeddedFieldEngineer)
			embeddedMappings = append(embeddedMappings, mappingFE)
		}

		// Uncomment if created_by is needed in _embedded field
		//embeddedCreatedBy := api.NewEmbeddedBasicUser(inc.CreatedUpdated.CreatedBy())
		//mapping := *hypermedia.EmbeddedResourcesMappingDefinition[embedded.CreatedBy].AddResource(embeddedCreatedBy)
		//embeddedMappings = append(embeddedMappings, mapping)

		incResp := api.IncidentResponse{
			Incident: p.convertIncidentToAPI(inc),
			Links:    p.resourceToHypermediaLinks(inc, hypermediaMapper, true),
			Embedded: p.resourceToEmbeddedField(inc, embeddedMappings, hypermediaMapper),
		}
		apiList = append(apiList, incResp)
	}

	pageInfo := api.PageInfo{
		Total: incidentList.Total,
		Size:  incidentList.Size,
		Page:  incidentList.Page,
	}

	resp := api.IncidentListResponse{
		Result:   apiList,
		PageInfo: pageInfo,
		Links:    p.hypermediaListLinks(hypermediaMapper, incidentList.Pagination),
	}

	p.renderJSON(w, resp)
}

func (p incidentPresenter) convertIncidentToAPI(inc incident.Incident) api.Incident {
	var timelogUUIDs []api.UUID
	for _, timelog := range inc.Timelogs {
		timelogUUIDs = append(timelogUUIDs, api.UUID(timelog))
	}

	var feUUID *api.UUID
	if inc.FieldEngineerID != nil {
		uuidS := api.UUID(inc.FieldEngineerID.String())
		feUUID = &uuidS
	}

	apiInc := api.Incident{
		UUID:             inc.UUID().String(),
		Number:           inc.Number,
		ExternalID:       inc.ExternalID,
		ShortDescription: inc.ShortDescription,
		Description:      inc.Description,
		FieldEngineer:    feUUID,
		State:            inc.State(),
		Timelogs:         timelogUUIDs,
		CreatedUpdated:   api.NewCreatedUpdatedInfo(inc.CreatedUpdated),
	}

	return apiInc
}
