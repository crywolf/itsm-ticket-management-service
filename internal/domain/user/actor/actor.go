package actor

import (
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user"
)

// Actor represents info about the user who initiated te API call
type Actor struct {
	BasicUser       user.BasicUser
	fieldEngineerID *ref.UUID
}

// ExternalUserUUID returns UUID of the actor in the external user microservice
func (e Actor) ExternalUserUUID() ref.ExternalUserUUID {
	return e.BasicUser.ExternalUserUUID
}

// IsFieldEngineer returns true if the actor is field engineer, otherwise it returns false
func (e Actor) IsFieldEngineer() bool {
	return e.fieldEngineerID != nil
}

// FieldEngineerID returns pointer to field engineer UUID if the actor is field engineer or nil if not
func (e Actor) FieldEngineerID() *ref.UUID {
	return e.fieldEngineerID
}

// SetFieldEngineerID sets field engineer UUID
func (e *Actor) SetFieldEngineerID(fieldEngineerID *ref.UUID) {
	if e.fieldEngineerID == nil {
		e.fieldEngineerID = fieldEngineerID
	}
}
