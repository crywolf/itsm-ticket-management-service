package mocks

import (
	"context"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/stretchr/testify/mock"
)

// IncidentServiceMock is an incident service mock
type IncidentServiceMock struct {
	mock.Mock
}

// CreateIncident mock
func (i *IncidentServiceMock) CreateIncident(ctx context.Context, channelID ref.ChannelID, params api.CreateIncidentParams) (ref.UUID, error) {
	panic("implement me")
}

// GetIncident mock
func (i *IncidentServiceMock) GetIncident(_ context.Context, channelID ref.ChannelID, ID ref.UUID) (incident.Incident, error) {
	args := i.Called(ID, channelID)
	return args.Get(0).(incident.Incident), args.Error(1)
}

// ListIncidents mock
func (i *IncidentServiceMock) ListIncidents(_ context.Context, channelID ref.ChannelID) ([]incident.Incident, error) {
	args := i.Called(channelID)
	return args.Get(0).([]incident.Incident), args.Error(1)
}
