package hypermedia

import (
	"net/url"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
)

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
// It must be implemented by domain object.
type ActionsMapper interface {
	// UUID return identifier of the domain object
	UUID() ref.UUID

	// AllowedActions returns list of actions that can be performed with the domain object according to its state and other conditions
	AllowedActions(actor actor.Actor) []string
}
