package tsession

import (
	"fmt"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
)

// TimeSession domain object
type TimeSession struct {
	uuid ref.UUID

	// State of the time session
	state State

	Incidents []IncidentInfo

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

// IncidentInfo contains basic info about incident in the session
type IncidentInfo struct {
	IncidentID         ref.UUID
	HasSupplierProduct bool
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

// StartWorking add incident to time session a sets it to Work state
func (e *TimeSession) StartWorking(inc incident.Incident) error {
	if e.state != StateTravel && e.state != StateWork {
		return domain.NewErrorf(domain.ErrorCodeUnknown, "time session is not in Travel nor Work state")
	}

	e.state = StateWork
	return e.AddIncident(inc)
}

// AddIncident adds incident to the time session (skips adding the same incident multiple times)
func (e *TimeSession) AddIncident(inc incident.Incident) error {
	//hasSp := inc.SupplierProduct != nil // TODO supplier product not implemented yet
	hasSp := true

	for _, info := range e.Incidents {
		// skip adding already added one
		if info.IncidentID == inc.UUID() {
			return nil
		}

		// prevent mixing tickets with and without supplier products
		if info.HasSupplierProduct != hasSp {
			return domain.NewErrorf(domain.ErrorCodeInvalidArgument, "cannot mix incidents with and without supplier product in the time session")
		}
	}

	if e.state != StateWork {
		return domain.NewErrorf(domain.ErrorCodeUnknown, "time session is not in Work state")
	}

	incInfo := IncidentInfo{
		IncidentID:         inc.UUID(),
		HasSupplierProduct: hasSp,
	}
	e.Incidents = append(e.Incidents, incInfo)

	return nil
}
