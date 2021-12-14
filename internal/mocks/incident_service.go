package mocks

import (
	"context"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	converters "github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/input_converters"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
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
	args := s.Called(ID, channelID, actor)
	return args.Get(0).(incident.Incident), args.Error(1)
}

// ListIncidents mock
func (s *IncidentServiceMock) ListIncidents(_ context.Context, channelID ref.ChannelID, actor actor.Actor, paginationParams converters.PaginationParams) (repository.IncidentList, error) {
	args := s.Called(channelID, actor, paginationParams)
	return args.Get(0).(repository.IncidentList), args.Error(1)
}
