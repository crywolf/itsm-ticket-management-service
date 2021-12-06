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
func (i *IncidentServiceMock) CreateIncident(_ context.Context, channelID ref.ChannelID, actor actor.Actor, params api.CreateIncidentParams) (ref.UUID, error) {
	args := i.Called(channelID, actor, params)
	return args.Get(0).(ref.UUID), args.Error(1)
}

// GetIncident mock
func (i *IncidentServiceMock) GetIncident(_ context.Context, channelID ref.ChannelID, actor actor.Actor, ID ref.UUID) (incident.Incident, error) {
	args := i.Called(ID, channelID, actor)
	return args.Get(0).(incident.Incident), args.Error(1)
}

// ListIncidents mock
func (i *IncidentServiceMock) ListIncidents(_ context.Context, channelID ref.ChannelID, actor actor.Actor, paginationParams converters.PaginationParams) (repository.IncidentList, error) {
	args := i.Called(channelID, actor, paginationParams)
	return args.Get(0).(repository.IncidentList), args.Error(1)
}
