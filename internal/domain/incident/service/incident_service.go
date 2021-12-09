package incidentsvc

import (
	"context"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	converters "github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/input_converters"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
)

// NewIncidentService creates the incident service
func NewIncidentService(r repository.IncidentRepository) IncidentService {
	return &incidentService{r}
}

type incidentService struct {
	r repository.IncidentRepository
}

func (s *incidentService) CreateIncident(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, params api.CreateIncidentParams) (ref.UUID, error) {
	// TODO validate if Caller or FE is not trying to set other fields then he is allowed to set, only SD agent can set everything (also in Update)
	newIncident := incident.Incident{
		ExternalID:       params.ExternalID,
		ShortDescription: params.ShortDescription,
		Description:      params.Description,
	}

	if err := newIncident.CreatedUpdated.SetCreatedBy(actor.BasicUser.UUID()); err != nil {
		return ref.UUID(""), err
	}
	if err := newIncident.CreatedUpdated.SetUpdatedBy(actor.BasicUser.UUID()); err != nil {
		return ref.UUID(""), err
	}

	return s.r.AddIncident(ctx, channelID, newIncident)
}

func (s *incidentService) GetIncident(ctx context.Context, channelID ref.ChannelID, _ actor.Actor, ID ref.UUID) (incident.Incident, error) {
	return s.r.GetIncident(ctx, channelID, ID)
}

func (s *incidentService) ListIncidents(ctx context.Context, channelID ref.ChannelID, _ actor.Actor, params converters.PaginationParams) (repository.IncidentList, error) {
	return s.r.ListIncidents(ctx, channelID, params.Page(), params.ItemsPerPage())
}
