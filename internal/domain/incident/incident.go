package incident

import (
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
)

// Incident domain object
type Incident struct {
	UUID ref.UUID

	// ID in external system
	ExternalID string

	ShortDescription string

	Description string

	state State

	CreatedUpdated types.CreatedUpdated
}

// GetState ...
func (i Incident) GetState() State {
	return i.state
}

// SetState ...
func (i *Incident) SetState(s State) error {
	i.state = s
	return nil
}
