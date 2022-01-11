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

// IncidentService provides incident operations
type IncidentService interface {
	// CreateIncident creates new incident and adds it to the repository
	CreateIncident(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, params api.CreateIncidentParams) (ref.UUID, error)

	// UpdateIncident updates the given incident in the repository
	UpdateIncident(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, ID ref.UUID, params api.UpdateIncidentParams) (ref.UUID, error)

	// GetIncident returns the incident with the given ID from the repository
	GetIncident(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, ID ref.UUID) (incident.Incident, error)

	// ListIncidents returns the list of incidents from the repository
	ListIncidents(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, paginationParams converters.PaginationParams) (repository.IncidentList, error)

	// StartWorking is used by actor (field engineer) to start working on the incident
	StartWorking(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, incID ref.UUID, params api.IncidentStartWorkingParams, clock domain.Clock) error

	// StopWorking is used by actor (field engineer) to stop working on the incident
	StopWorking(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, incID ref.UUID, params api.IncidentStopWorkingParams, clock domain.Clock) error
}
