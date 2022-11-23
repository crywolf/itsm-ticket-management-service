package hypermedia

import (
	"context"
	"net/url"

	"github.com/crywolf/itsm-ticket-management-service/internal/domain/embedded"
	fieldengineersvc "github.com/crywolf/itsm-ticket-management-service/internal/domain/field_engineer/service"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user/actor"
)

type IncidentMapper interface {
	Mapper
	FieldEngineerSvc() fieldengineersvc.FieldEngineerService
	Ctx() context.Context
	ChannelID() ref.ChannelID
}

// Mapper maps domain object to hypermedia representation
type Mapper interface {
	// RoutesToHypermediaActionLinks maps domain object actions to hypermedia action links.
	// It must be implemented for specific resource in object that includes BaseHypermediaMapper object via composition.
	RoutesToHypermediaActionLinks() ActionLinks

	// SelfLink returns 'self' link URL. It is automatically implemented via BaseHypermediaMapper object.
	SelfLink() string

	// RequestURL returns current URL
	RequestURL() *url.URL

	// ServerAddr returns server URL. It is automatically implemented via BaseHypermediaMapper object.
	ServerAddr() string

	// Actor returns user who initiated current API call
	Actor() actor.Actor
}

// ActionsMapper provides domain object mapping to hypermedia actions.
// It must be implemented by the domain object.
type ActionsMapper interface {
	// UUID returns identifier of the domain object
	UUID() ref.UUID

	// AllowedActions returns list of actions that can be performed with the domain object according to its state and other conditions
	AllowedActions(actor actor.Actor) []string
}

// EmbeddedResourceMapper provides information about domain object mapping to hypermedia embedded resources.
// It must be implemented by the domain object.
type EmbeddedResourceMapper interface {
	// EmbeddedResources returns list of other resources that are 'embedded' in the resource
	EmbeddedResources(actor actor.Actor) []embedded.Resource
}

// EmbeddedResource can append 'self' link to its '_links' object
type EmbeddedResource interface {
	// UUID returns resource's UUID
	UUID() string

	// AppendSelfLink adds resource's 'self' link to its '_links' object
	AppendSelfLink(url string)
}
