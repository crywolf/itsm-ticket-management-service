package timelog

import (
	"errors"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
)

// Timelog domain object
type Timelog struct {
	uuid ref.UUID

	Remote bool

	Work uint

	VisitSummary string

	CreatedUpdated types.CreatedUpdated
}

// UUID getter
func (i Timelog) UUID() ref.UUID {
	return i.uuid
}

// SetUUID returns error if UUID was already set
func (i *Timelog) SetUUID(v ref.UUID) error {
	if !i.uuid.IsZero() {
		return errors.New("cannot set UUID, it was already set")
	}
	i.uuid = v
	return nil
}

