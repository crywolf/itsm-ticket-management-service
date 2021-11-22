package incident

import (
	"errors"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident/timelog"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
)

// Incident domain object
type Incident struct {
	uuid ref.UUID

	Number string

	// ID in external system
	ExternalID string

	ShortDescription string

	Description string

	state State

	openTimelog *timelog.Timelog

	Timelogs []ref.UUID

	CreatedUpdated types.CreatedUpdated
}

// UUID getter
func (i Incident) UUID() ref.UUID {
	return i.uuid
}

// SetUUID returns error if UUID was already set
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

func (i Incident) OpenTimelog() *timelog.Timelog {
	return i.openTimelog
}

func (i *Incident) SetOpenTimelog(openTimelog *timelog.Timelog) {
	i.openTimelog = openTimelog
}

func (i Incident) HasOpenTimelog() bool {
	return i.openTimelog != nil
}

// AllowedAction represents action that can be performed with the incident
type AllowedAction string

func (a AllowedAction) String() string {
	return string(a)
}

// AllowedActions values
const (
	ActionCancel       AllowedAction = "Cancel"
	ActionStartWorking AllowedAction = "StartWorking"
)

// AllowedActions returns list of actions that can be performed with the incident according to its state and other conditions
func (i Incident) AllowedActions() []string {
	// TODO - this is just for testing
	var acts []string
	acts = append(acts, ActionCancel.String(), ActionStartWorking.String())
	return acts
}
