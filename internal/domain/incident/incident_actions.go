package incident

import (
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
)

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
func (e Incident) AllowedActions(actor actor.Actor) []string {
	var acts []string
	if err := e.canBeCancelled(actor); err == nil {
		acts = append(acts, ActionCancel.String())
	}

	if err := e.canStartWorking(actor); err == nil {
		acts = append(acts, ActionStartWorking.String())
	}

	return acts
}

// Cancel cancels the ticket
func (e *Incident) Cancel(actor actor.Actor) error {
	if err := e.canBeCancelled(actor); err != nil {
		return err
	}

	return nil
}

func (e *Incident) canBeCancelled(actor actor.Actor) error {
	if e.state != StateNew {
		return domain.NewErrorf(domain.ErrorCodeInvalidArgument, "ticket can be cancelled only in NEW state")
	}

	return nil
}

// StartWorking can be used by assigned field engineer to start working on the ticket
func (e *Incident) StartWorking(actor actor.Actor) error {
	if err := e.canStartWorking(actor); err != nil {
		return err
	}

	if err := e.SetState(StateInProgress); err != nil {
		return err
	}
	// TODO - open timelog ...
	return nil
}

func (e *Incident) canStartWorking(actor actor.Actor) error {
	if !actor.IsFieldEngineer() {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "user is not field engineer, only assigned field engineer can start working")
	}

	if e.fieldEngineer == nil {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "ticket does not have any field engineer assigned")
	}

	if actor.IsFieldEngineer() && e.fieldEngineer != nil && actor.FieldEngineer().UUID() != e.fieldEngineer.UUID() {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "user is not an assigned field engineer, only assigned field engineer can start working")
	}

	if e.HasOpenTimelog() {
		return domain.NewErrorf(domain.ErrorCodeInvalidArgument, "ticket already has an open timelog")
	}

	return nil
}
