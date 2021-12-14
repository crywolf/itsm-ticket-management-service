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
func NewIncidentService(repo repository.IncidentRepository) IncidentService {
	return &incidentService{repo}
}

type incidentService struct {
	repo repository.IncidentRepository
}

func (s *incidentService) CreateIncident(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, params api.CreateIncidentParams) (ref.UUID, error) {
	// TODO validate that Caller or FE is not trying to set other fields then he is allowed to set, only SD agent can set everything (also in Update)
	newIncident := incident.Incident{
		Number:           params.Number,
		ExternalID:       params.ExternalID,
		ShortDescription: params.ShortDescription,
		Description:      params.Description,
	}

	if err := newIncident.CreatedUpdated.SetCreatedBy(actor.BasicUser); err != nil {
		return ref.UUID(""), err
	}
	if err := newIncident.CreatedUpdated.SetUpdatedBy(actor.BasicUser); err != nil {
		return ref.UUID(""), err
	}
	return s.repo.AddIncident(ctx, channelID, newIncident)
}

// UpdateIncident updates the given incident in the repository
func (s *incidentService) UpdateIncident(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, ID ref.UUID, params api.UpdateIncidentParams) (ref.UUID, error) {
	inc, err := s.repo.GetIncident(ctx, channelID, ID)
	if err != nil {
		return ref.UUID(""), err
	}

	inc.ShortDescription = params.ShortDescription
	inc.Description = params.Description

	if err := inc.CreatedUpdated.SetUpdatedBy(actor.BasicUser); err != nil {
		return ref.UUID(""), err
	}

	return s.repo.UpdateIncident(ctx, channelID, inc)
}

func (s *incidentService) GetIncident(ctx context.Context, channelID ref.ChannelID, _ actor.Actor, ID ref.UUID) (incident.Incident, error) {
	return s.repo.GetIncident(ctx, channelID, ID)
}

func (s *incidentService) ListIncidents(ctx context.Context, channelID ref.ChannelID, _ actor.Actor, params converters.PaginationParams) (repository.IncidentList, error) {
	return s.repo.ListIncidents(ctx, channelID, params.Page(), params.ItemsPerPage())
}
