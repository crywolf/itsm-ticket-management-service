package fieldengineer

import (
	"errors"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
)

// FieldEngineer represents info about the user who initiated te API call
type FieldEngineer struct {
	uuid ref.UUID

	BasicUserUUID ref.UUID
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
