package hypermedia

import (
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
)

// ActionLinks maps allowed resource domain actions to hypermedia action links
type ActionLinks map[string]api.ActionLink

// Mapper maps domain object to hypermedia representation
type Mapper interface {
	// RoutesToHypermediaActionLinks maps domain object actions to hypermedia action links.
	// It must be implemented for specific resource in object that includes BaseHypermediaMapper object via composition.
	RoutesToHypermediaActionLinks() ActionLinks
	// SelfLink returns 'self' link URL. It is automatically implemented via BaseHypermediaMapper object.
	SelfLink() string
	// ServerAddr returns server URL. It is automatically implemented via BaseHypermediaMapper object.
	ServerAddr() string
}

// ActionsMapper provides domain object mapping to hypermedia actions.
// It must be implemented by domain object.
type ActionsMapper interface {
	AllowedActions() []string
	UUID() ref.UUID
}
