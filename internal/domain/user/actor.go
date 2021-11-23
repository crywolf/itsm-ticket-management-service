package user

import (
	fieldengineer "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
)

// Actor represents info about the user who initiated te API call
type Actor struct {
	BasicUser     BasicUser
	fieldEngineer *fieldengineer.FieldEngineer
}

// ExternalUserUUID returns UUID of the actor in the external user microservice
func (e Actor) ExternalUserUUID() ref.ExternalUserUUID {
	return e.BasicUser.ExternalUserUUID
}

// IsFieldEngineer returns true if the actor is field engineer, otherwise it returns false
func (e Actor) IsFieldEngineer() bool {
	return e.fieldEngineer != nil
}

// FieldEngineer returns pointer to field engineer if the actor is field engineer or nil if not
func (e Actor) FieldEngineer() *fieldengineer.FieldEngineer {
	return e.fieldEngineer
}

// SetFieldEngineer sets field engineer
func (e *Actor) SetFieldEngineer(fieldEngineer *fieldengineer.FieldEngineer) {
	if e.fieldEngineer == nil {
		e.fieldEngineer = fieldEngineer
	}
}
