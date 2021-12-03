package incidentsvc

import (
	"context"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
)

// NewIncidentService creates the incident service
func NewIncidentService(r repository.IncidentRepository) IncidentService {
	return &service{r}
}

type service struct {
	r repository.IncidentRepository
}

func (s *service) CreateIncident(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, params api.CreateIncidentParams) (ref.UUID, error) {
	newIncident := incident.Incident{
		ExternalID:       params.ExternalID,
		ShortDescription: params.ShortDescription,
		Description:      params.Description,
	}

	if err := newIncident.CreatedUpdated.SetCreatedBy(actor.ExternalUserUUID()); err != nil {
		return ref.UUID(""), err
	}
	if err := newIncident.CreatedUpdated.SetUpdatedBy(actor.ExternalUserUUID()); err != nil {
		return ref.UUID(""), err
	}

	return s.r.AddIncident(ctx, channelID, newIncident)
}

func (s *service) GetIncident(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, ID ref.UUID) (incident.Incident, error) {
	return s.r.GetIncident(ctx, channelID, ID)
}

func (s *service) ListIncidents(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, page, perPage uint) (repository.IncidentList, error) {
	return s.r.ListIncidents(ctx, channelID, page, perPage)
}
