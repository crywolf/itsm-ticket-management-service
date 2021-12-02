package repository

import (
	"context"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
)

// BasicUserRepository provides access to the Basic User repository
type BasicUserRepository interface {
	// GetBasicUser returns the Basic User with the given ID from the repository
	GetBasicUser(ctx context.Context, channelID ref.ChannelID, ID ref.UUID) (user.BasicUser, error)

	// GetBasicUserByExternalID returns the Basic User with the given external ID from the repository
	GetBasicUserByExternalID(ctx context.Context, channelID ref.ChannelID, externalID ref.ExternalUserUUID) (user.BasicUser, error)
}

// IncidentRepository provides access to the incidents repository
type IncidentRepository interface {
	// AddIncident adds the given incident to the repository
	AddIncident(ctx context.Context, channelID ref.ChannelID, inc incident.Incident) (ref.UUID, error)

	// GetIncident returns the incident with the given ID from the repository
	GetIncident(ctx context.Context, channelID ref.ChannelID, ID ref.UUID) (incident.Incident, error)

	// ListIncidents returns the list of incidents from the repository
	ListIncidents(ctx context.Context, channelID ref.ChannelID) ([]incident.Incident, error)
}
