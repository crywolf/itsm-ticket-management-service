package incident

import (
	"fmt"

	"github.com/crywolf/itsm-ticket-management-service/internal/domain"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/embedded"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/incident/timelog"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/types"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user/actor"
)

// Incident domain object
type Incident struct {
	uuid ref.UUID

	Number string

	// ID in external system
	ExternalID string

	ShortDescription string

	Description string

	FieldEngineerID *ref.UUID

	state State

	openTimelog *timelog.Timelog

	// TODO make it private - timelogIDs
	Timelogs []ref.UUID

	CreatedUpdated types.CreatedUpdated
}

// New creates initialized Incident
func New(clock domain.Clock, basicUser user.BasicUser) (Incident, error) {
	inc := Incident{state: StateNew}

	now := clock.NowFormatted()

	if err := inc.CreatedUpdated.SetCreated(basicUser, now); err != nil {
		return inc, err
	}

	if err := inc.CreatedUpdated.SetUpdated(basicUser, now); err != nil {
		return inc, err
	}

	return inc, nil
}

// UUID getter
func (e Incident) UUID() ref.UUID {
	return e.uuid
}

// SetUUID returns error if UUID was already set
func (e *Incident) SetUUID(v ref.UUID) error {
	if !e.uuid.IsZero() {
		return fmt.Errorf("incident: cannot set UUID, it was already set (%s)", e.uuid)
	}
	e.uuid = v
	return nil
}

// State getter
func (e Incident) State() State {
	return e.state
}

// SetState ...
func (e *Incident) SetState(s State) error {
	// TODO add state machine and checks
	e.state = s
	return nil
}

// OpenTimelog returns open timelog if any or nil pointer
func (e Incident) OpenTimelog() *timelog.Timelog {
	return e.openTimelog
}

// SetOpenTimelog sets open timelog (do not use in the domain, method is used by repository)
func (e *Incident) SetOpenTimelog(openTimelog *timelog.Timelog) {
	e.openTimelog = openTimelog
}

// HasOpenTimelog returns true if the ticket has an open timelog
func (e Incident) HasOpenTimelog() bool {
	return e.openTimelog != nil
}

// EmbeddedResources returns list of other objects that are 'embedded' in the ticket
func (e Incident) EmbeddedResources(actor actor.Actor) []embedded.Resource {
	var resources []embedded.Resource

	resources = append(resources, embedded.FieldEngineer)
	resources = append(resources, e.CreatedUpdated.EmbeddedResources()...)

	// TODO add other fields...
	if actor.IsFieldEngineer() {
	}
	return resources
}
