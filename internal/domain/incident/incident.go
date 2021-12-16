package incident

import (
	"errors"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/embedded"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident/timelog"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
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

// UUID getter
func (e Incident) UUID() ref.UUID {
	return e.uuid
}

// SetUUID returns error if UUID was already set
func (e *Incident) SetUUID(v ref.UUID) error {
	if !e.uuid.IsZero() {
		return errors.New("incident: cannot set UUID, it was already set")
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

// SetOpenTimelog sets open timelog (do not use in the domain,method is used by repository)
func (e *Incident) SetOpenTimelog(openTimelog *timelog.Timelog) {
	e.openTimelog = openTimelog
}

// HasOpenTimelog returns true if the ticket has open timelog
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
