package mocks

import (
	"context"

	"github.com/crywolf/itsm-ticket-management-service/internal/domain"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/incident"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/incident/timelog"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/crywolf/itsm-ticket-management-service/internal/http/rest/api"
	converters "github.com/crywolf/itsm-ticket-management-service/internal/http/rest/input_converters"
	"github.com/crywolf/itsm-ticket-management-service/internal/repository"
	"github.com/stretchr/testify/mock"
)

// IncidentServiceMock is an incident service mock
type IncidentServiceMock struct {
	mock.Mock
}

// CreateIncident mock
func (s *IncidentServiceMock) CreateIncident(_ context.Context, channelID ref.ChannelID, actor actor.Actor, params api.CreateIncidentParams) (ref.UUID, error) {
	args := s.Called(channelID, actor, params)
	return args.Get(0).(ref.UUID), args.Error(1)
}

// UpdateIncident updates the given incident in the repository
func (s *IncidentServiceMock) UpdateIncident(_ context.Context, channelID ref.ChannelID, actor actor.Actor, ID ref.UUID, params api.UpdateIncidentParams) (ref.UUID, error) {
	args := s.Called(channelID, actor, ID, params)
	return args.Get(0).(ref.UUID), args.Error(1)
}

// GetIncident mock
func (s *IncidentServiceMock) GetIncident(_ context.Context, channelID ref.ChannelID, actor actor.Actor, ID ref.UUID) (incident.Incident, error) {
	args := s.Called(channelID, actor, ID)
	return args.Get(0).(incident.Incident), args.Error(1)
}

// ListIncidents mock
func (s *IncidentServiceMock) ListIncidents(_ context.Context, channelID ref.ChannelID, actor actor.Actor, paginationParams converters.PaginationParams) (repository.IncidentList, error) {
	args := s.Called(channelID, actor, paginationParams)
	return args.Get(0).(repository.IncidentList), args.Error(1)
}

// StartWorking mock
func (s *IncidentServiceMock) StartWorking(_ context.Context, channelID ref.ChannelID, actor actor.Actor, incID ref.UUID, params api.IncidentStartWorkingParams, _ domain.Clock) error {
	args := s.Called(channelID, actor, incID, params)
	return args.Error(0)
}

// StopWorking mock
func (s *IncidentServiceMock) StopWorking(_ context.Context, channelID ref.ChannelID, actor actor.Actor, incID ref.UUID, params api.IncidentStopWorkingParams, _ domain.Clock) error {
	args := s.Called(channelID, actor, incID, params)
	return args.Error(0)
}

// GetIncidentTimelog mock
func (s *IncidentServiceMock) GetIncidentTimelog(_ context.Context, channelID ref.ChannelID, actor actor.Actor, incID ref.UUID, timelogID ref.UUID) (timelog.Timelog, error) {
	args := s.Called(channelID, actor, incID, timelogID)
	return args.Get(0).(timelog.Timelog), args.Error(1)
}
