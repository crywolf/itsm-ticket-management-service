package incident

import (
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident/timelog"
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

	if err := e.SetState(StateCancelled); err != nil {
		return err
	}

	return nil
}

func (e *Incident) canBeCancelled(_ actor.Actor) error {
	if e.state != StateNew {
		return domain.NewErrorf(domain.ErrorCodeInvalidArgument, "ticket can be cancelled only in NEW state")
	}

	return nil
}

// StartWorking can be used by assigned field engineer to start working on the ticket
func (e *Incident) StartWorking(actor actor.Actor, remote bool) error {
	if err := e.canStartWorking(actor); err != nil {
		return err
	}

	// open new timelog
	newTimelog := &timelog.Timelog{
		Remote: remote,
	}
	if err := newTimelog.CreatedUpdated.SetCreatedBy(actor.BasicUser); err != nil {
		return err
	}
	if err := newTimelog.CreatedUpdated.SetUpdatedBy(actor.BasicUser); err != nil {
		return err
	}

	e.openTimelog = newTimelog

	if err := e.SetState(StateInProgress); err != nil {
		return err
	}

	return nil
}

func (e *Incident) canStartWorking(actor actor.Actor) error {
	if !actor.IsFieldEngineer() {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "user is not field engineer, only assigned field engineer can start working")
	}

	if e.FieldEngineerID == nil {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "ticket does not have any field engineer assigned")
	}

	if actor.IsFieldEngineer() && e.FieldEngineerID != nil && *actor.FieldEngineerID() != *e.FieldEngineerID {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "user is not assigned as field engineer, only assigned field engineer can start working")
	}

	// TODO disallow if FE did not accepted the incident

	if e.HasOpenTimelog() {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "ticket already has an open timelog")
	}

	if e.state != StateNew && e.state != StateInProgress && e.state != StateOnHold {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "ticket is not in New, InProgress nor OnHold state")
	}

	return nil
}
