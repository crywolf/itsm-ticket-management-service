package api

import (
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
)

// BasicUser API object
// swagger:model
type BasicUser struct {
	// required: true
	// swagger:strfmt uuid
	UUID string `json:"uuid"`

	// User in external microservice
	ExternalUserUUID ref.ExternalUserUUID `json:"external_user_uuid,omitempty"`

	// required: true
	Name string `json:"name"`

	// required: true
	Surname string `json:"surname"`

	// example: KompiTech
	OrgDisplayName string `json:"org_display_name,omitempty"`

	// example: a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com
	OrgName string `json:"org_name,omitempty"`

	CreatedUpdated
}
