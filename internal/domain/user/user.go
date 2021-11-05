package user

import "github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"

// BasicUserUUID is a reference to a Basic User resource
type BasicUserUUID ref.UUID

// BasicUser represents basic info about the user in external service (user microservice)
type BasicUser struct {
	UUID BasicUserUUID

	// User in external microservice
	ExternalUserUUID ref.ExternalUserUUID

	Name string

	Surname string

	// example: KompiTech
	OrgDisplayName string

	// example: a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com
	OrgName string
}
