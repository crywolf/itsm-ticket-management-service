package incidentsvc

import (
	"context"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
)

// NewIncidentService creates an incident service
func NewIncidentService(r IncidentRepository) IncidentService {
	return &service{r}
}

type service struct {
	r IncidentRepository
}

func (s *service) CreateIncident(ctx context.Context, channelID ref.ChannelID, actor user.Actor, params api.CreateIncidentParams) (ref.UUID, error) {
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

func (s *service) GetIncident(ctx context.Context, channelID ref.ChannelID, actor user.Actor, ID ref.UUID) (incident.Incident, error) {
	return s.r.GetIncident(ctx, channelID, ID)
}

func (s *service) ListIncidents(ctx context.Context, channelID ref.ChannelID, actor user.Actor) ([]incident.Incident, error) {
	return s.r.ListIncidents(ctx, channelID)
}
