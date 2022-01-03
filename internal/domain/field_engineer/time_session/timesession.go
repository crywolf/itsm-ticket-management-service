package tsession

import (
	"fmt"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
)

// TimeSession domain object
type TimeSession struct {
	uuid ref.UUID

	// State of the time session
	state State

	Incidents []ref.UUID

	// Time spent working (in seconds) counted from all timespans
	Work uint

	// Time spent travelling to customer (in seconds)
	Travel uint

	// Time spent travelling from customer (in seconds)
	TravelBack uint

	// Distance travelled to the customer and back in 'travel units'
	TravelDistanceInTravelUnits uint

	CreatedUpdated types.CreatedUpdated
}

// UUID getter
func (e TimeSession) UUID() ref.UUID {
	return e.uuid
}

// SetUUID returns error if UUID was already set
func (e *TimeSession) SetUUID(v ref.UUID) error {
	if !e.uuid.IsZero() {
		return fmt.Errorf("time session: cannot set UUID, it was already set (%s)", e.uuid)
	}
	e.uuid = v
	return nil
}

// State getter
func (e TimeSession) State() State {
	return e.state
}

// SetState ...
func (e *TimeSession) SetState(s State) error {
	// TODO add state machine and checks
	e.state = s
	return nil
}

// AddIncident adds incident to time session
func (e *TimeSession) AddIncident(incidentID ref.UUID) error {
	return nil
}
