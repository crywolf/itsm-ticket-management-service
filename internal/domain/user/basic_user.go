package user

import (
	"errors"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
)

// DefaultItemsPerPage is a default number of records to be shown per one page in the list of items
const DefaultItemsPerPage uint = 10

// BasicUser represents basic info about the user in external service (user microservice)
type BasicUser struct {
	uuid ref.UUID

	// User in external microservice
	ExternalUserUUID ref.ExternalUserUUID

	Name string

	Surname string

	// example: KompiTech
	OrgDisplayName string

	// example: a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com
	OrgName string
}

// UUID getter
func (e BasicUser) UUID() ref.UUID {
	return e.uuid
}

// SetUUID returns error if UUID was already set
func (e *BasicUser) SetUUID(v ref.UUID) error {
	if !e.uuid.IsZero() {
		return errors.New("basic user: cannot set UUID, it was already set")
	}
	e.uuid = v
	return nil
}

// ItemsPerPage returns a number of records to be shown per one page in the list of items for this user
func (e BasicUser) ItemsPerPage() uint {
	return DefaultItemsPerPage
}
