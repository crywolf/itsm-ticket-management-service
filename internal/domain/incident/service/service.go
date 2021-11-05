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

func (s *service) CreateIncident(ctx context.Context, channelID ref.ChannelID /*, actor user.BasicUser,*/, params api.CreateIncidentParams) (ref.UUID, error) {
	newIncident := incident.Incident{
		ExternalID:       params.ExternalID,
		ShortDescription: params.ShortDescription,
		Description:      params.Description,
	}

	// TODO take from func params
	actor := user.BasicUser{
		ExternalUserUUID: "8183eaca-56c0-41d9-9291-1d295dd53763",
	}
	err := newIncident.CreatedUpdated.SetCreatedBy(actor.ExternalUserUUID)
	if err != nil {
		return ref.UUID(""), err
	}

	return s.r.AddIncident(ctx, channelID, newIncident)
}

func (s *service) GetIncident(ctx context.Context, channelID ref.ChannelID, ID ref.UUID) (incident.Incident, error) {
	return s.r.GetIncident(ctx, channelID, ID)
}

func (s *service) ListIncidents(ctx context.Context, channelID ref.ChannelID) ([]incident.Incident, error) {
	return s.r.ListIncidents(ctx, channelID)
}
