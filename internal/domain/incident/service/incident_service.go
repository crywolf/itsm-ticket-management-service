package incidentsvc

import (
	"context"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident/timelog"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	converters "github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/input_converters"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
)

// NewIncidentService creates the incident service
func NewIncidentService(incidentRepository repository.IncidentRepository, fieldEngineerRepository repository.FieldEngineerRepository) IncidentService {
	return &incidentService{
		incidentRepository:      incidentRepository,
		fieldEngineerRepository: fieldEngineerRepository,
	}
}

type incidentService struct {
	incidentRepository      repository.IncidentRepository
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
	if err := newIncident.SetState(incident.StateNew); err != nil {
		return ref.UUID(""), err
	}

	if err := newIncident.CreatedUpdated.SetCreatedBy(actor.BasicUser); err != nil {
		return ref.UUID(""), err
	}
	if err := newIncident.CreatedUpdated.SetUpdatedBy(actor.BasicUser); err != nil {
		return ref.UUID(""), err
	}

	return s.incidentRepository.AddIncident(ctx, channelID, newIncident)
}

// UpdateIncident updates the given incident in the repository
func (s *incidentService) UpdateIncident(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, ID ref.UUID, params api.UpdateIncidentParams) (ref.UUID, error) {
	inc, err := s.incidentRepository.GetIncident(ctx, channelID, ID)
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

	return s.incidentRepository.UpdateIncident(ctx, channelID, inc)
}

func (s *incidentService) GetIncident(ctx context.Context, channelID ref.ChannelID, _ actor.Actor, ID ref.UUID) (incident.Incident, error) {
	return s.incidentRepository.GetIncident(ctx, channelID, ID)
}

func (s *incidentService) ListIncidents(ctx context.Context, channelID ref.ChannelID, _ actor.Actor, params converters.PaginationParams) (repository.IncidentList, error) {
	return s.incidentRepository.ListIncidents(ctx, channelID, params.Page(), params.ItemsPerPage())
}

func (s *incidentService) StartWorking(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, incID ref.UUID, params api.IncidentStartWorkingParams, clock domain.Clock) error {
	if !actor.IsFieldEngineer() {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "actor is not field engineer")
	}
	feID := actor.FieldEngineerID()

	inc, err := s.incidentRepository.GetIncident(ctx, channelID, incID)
	if err != nil {
		return err
	}

	fe, err := s.fieldEngineerRepository.GetFieldEngineer(ctx, channelID, *feID)
	if err != nil {
		return err
	}

	if err := inc.StartWorking(actor, clock, params.Remote); err != nil {
		return err
	}

	if err := fe.StartWorking(actor, inc); err != nil {
		return err
	}

	if _, err := s.incidentRepository.UpdateIncident(ctx, channelID, inc); err != nil {
		return err
	}

	if _, err := s.fieldEngineerRepository.UpdateFieldEngineer(ctx, channelID, fe); err != nil {
		return err
	}

	return nil
}

func (s *incidentService) StopWorking(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, incID ref.UUID, params api.IncidentStopWorkingParams, clock domain.Clock) error {
	inc, err := s.incidentRepository.GetIncident(ctx, channelID, incID)
	if err != nil {
		return err
	}

	if err := inc.StopWorking(actor, clock, params.VisitSummary); err != nil {
		return err
	}

	if _, err := s.incidentRepository.UpdateIncident(ctx, channelID, inc); err != nil {
		return err
	}

	return nil
}

// GetIncidentTimelog returns the incident's timelog with the given ID from the repository
func (s *incidentService) GetIncidentTimelog(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, incID ref.UUID, timelogID ref.UUID) (timelog.Timelog, error) {
	return s.incidentRepository.GetIncidentTimelog(ctx, channelID, incID, timelogID)
}
