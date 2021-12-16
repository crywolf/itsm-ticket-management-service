package incidentsvc

import (
	"context"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	converters "github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/input_converters"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
)

// NewIncidentService creates the incident service
func NewIncidentService(repo repository.IncidentRepository, fieldEngineerRepository repository.FieldEngineerRepository) IncidentService {
	return &incidentService{
		repo:                    repo,
		fieldEngineerRepository: fieldEngineerRepository,
	}
}

type incidentService struct {
	repo                    repository.IncidentRepository
	fieldEngineerRepository repository.FieldEngineerRepository
}

func (s *incidentService) CreateIncident(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, params api.CreateIncidentParams) (ref.UUID, error) {
	// TODO validate that Caller or FE is not trying to set other fields then he is allowed to set, only SD agent can set everything (also in Update)

	var feUUID ref.UUID
	if params.FieldEngineerID != nil {
		feUUID = ref.UUID(*params.FieldEngineerID)
		if _, err := s.fieldEngineerRepository.GetFieldEngineer(ctx, channelID, feUUID); err != nil {
			return ref.UUID(""), domain.WrapErrorf(err, domain.ErrorCodeNotFound, "cannot assign field engineer")
		}
	}

	newIncident := incident.Incident{
		Number:           params.Number,
		ExternalID:       params.ExternalID,
		ShortDescription: params.ShortDescription,
		Description:      params.Description,
		FieldEngineerID:  &feUUID,
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

	if params.FieldEngineerID != nil {
		feUUID := ref.UUID(*params.FieldEngineerID)
		if _, err := s.fieldEngineerRepository.GetFieldEngineer(ctx, channelID, feUUID); err != nil {
			return ref.UUID(""), domain.WrapErrorf(err, domain.ErrorCodeNotFound, "cannot assign field engineer")
		}

		inc.FieldEngineerID = &feUUID
	}

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
