package timelog

import (
	"fmt"

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
func (e Timelog) UUID() ref.UUID {
	return e.uuid
}

// SetUUID returns error if UUID was already set
func (e *Timelog) SetUUID(v ref.UUID) error {
	if !e.uuid.IsZero() {
		return fmt.Errorf("timelog: cannot set UUID, it was already set (%s)", e.uuid)
	}
	e.uuid = v
	return nil
}
