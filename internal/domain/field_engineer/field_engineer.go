package fieldengineer

import (
	"errors"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
)

// FieldEngineer represents a user who can solve the tickets
type FieldEngineer struct {
	uuid ref.UUID

	BasicUser user.BasicUser

	CreatedUpdated types.CreatedUpdated
}

// UUID getter
func (e FieldEngineer) UUID() ref.UUID {
	return e.uuid
}

// SetUUID returns error if UUID was already set
func (e *FieldEngineer) SetUUID(v ref.UUID) error {
	if !e.uuid.IsZero() {
		return errors.New("field engineer: cannot set UUID, it was already set")
	}
	e.uuid = v
	return nil
}
