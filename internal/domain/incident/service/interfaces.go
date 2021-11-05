package incidentsvc

import (
	"context"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
)

// IncidentService provides incident operations
type IncidentService interface {
	// CreateIncident creates new incident and adds it to the repository
	CreateIncident(ctx context.Context, channelID ref.ChannelID, /* actor user.BasicUser,*/ params api.CreateIncidentParams) (ref.UUID, error)

	// GetIncident returns the incident with given ID from the repository
	GetIncident(ctx context.Context, channelID ref.ChannelID, ID ref.UUID) (incident.Incident, error)

	// ListIncidents returns the incidents from the repository
	ListIncidents(ctx context.Context, channelID ref.ChannelID) ([]incident.Incident, error)
}

// IncidentRepository provides reading access to the incidents repository
type IncidentRepository interface {
	// AddIncident adds the given incident to the repository
	AddIncident(ctx context.Context, channelID ref.ChannelID, inc incident.Incident) (ref.UUID, error)

	// GetIncident returns the incident with given ID from the repository
	GetIncident(ctx context.Context, channelID ref.ChannelID, ID ref.UUID) (incident.Incident, error)

	// ListIncidents returns the incidents from the repository
	ListIncidents(ctx context.Context, channelID ref.ChannelID) ([]incident.Incident, error)
}
