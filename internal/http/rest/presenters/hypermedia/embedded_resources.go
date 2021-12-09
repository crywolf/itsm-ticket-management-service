package hypermedia

import (
	"strings"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/embedded"
)

// TODO zabudovat do base presenteru

// EmbeddedResourcesMappingDefinition maps domain resource names to its hypermedia (embedded) versions
var EmbeddedResourcesMappingDefinition = map[embedded.Resource]*EmbeddedResourceMapping{
	embedded.CreatedBy: {
		ResourceName: embedded.CreatedBy,
		Key:          "created_by",
		Route:        "/basic_users/{uuid}",
	},
	embedded.UpdatedBy: {
		ResourceName: embedded.UpdatedBy,
		Key:          "updated_by",
		Route:        "/basic_users/{uuid}",
	},
}

// EmbeddedResourceMapping maps domain resource name to its hypermedia (embedded) version
type EmbeddedResourceMapping struct {
	ResourceName embedded.Resource
	Key          string
	Route        string
	Resource     EmbeddedResource
}

// AddResource ads new EmbeddedResource to the mapping
func (m *EmbeddedResourceMapping) AddResource(resource EmbeddedResource) *EmbeddedResourceMapping {
	m.Route = strings.ReplaceAll(m.Route, "{uuid}", resource.UUID())
	m.Resource = resource
	return m
}
