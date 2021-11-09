package incident

import (
	"errors"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
)

// Incident domain object
type Incident struct {
	uuid ref.UUID

	// ID in external system
	ExternalID string

	ShortDescription string

	Description string

	state State

	CreatedUpdated types.CreatedUpdated
}

// UUID getter
func (i Incident) UUID() ref.UUID {
	return i.uuid
}

// SetUUID ...
func (i *Incident) SetUUID(v ref.UUID) error {
	if !i.uuid.IsZero() {
		return errors.New("cannot set UUID, it was already set")
	}
	i.uuid = v
	return nil
}

// State getter
func (i Incident) State() State {
	return i.state
}

// SetState ...
func (i *Incident) SetState(s State) error {
	// TODO add state machine and checks
	i.state = s
	return nil
}
