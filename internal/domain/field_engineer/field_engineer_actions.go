package fieldengineer

import (
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	tsession "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer/time_session"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
)

// AllowedAction represents action that can be performed with the incident
type AllowedAction string

func (a AllowedAction) String() string {
	return string(a)
}

// AllowedActions values
const (
	ActionStartTravelling     AllowedAction = "StartTravelling"
	ActionStartTravellingBack AllowedAction = "StartTravellingBack"
	ActionCloseTimeSession    AllowedAction = "CloseTimeSession"
	//	ActionStartWorking    AllowedAction = "StartWorking"
	//  ActionStopWorking AllowedAction = "StopWorking"
)

// AllowedActions returns list of actions that can be performed with the time session according to its state and other conditions
func (e FieldEngineer) AllowedActions(actor actor.Actor) []string {
	var acts []string

	if err := e.canStartTravelling(actor); err == nil {
		acts = append(acts, ActionStartTravelling.String())
	}

	return acts
}

// StartWorking can be used by assigned field engineer to start working on the ticket
func (e *FieldEngineer) StartWorking(actor actor.Actor, inc incident.Incident) error {
	if err := e.canStartWorking(actor); err != nil {
		return err
	}

	if !e.HasOpenTimeSession() {
		// open new time session
		newTimeSession := &tsession.TimeSession{}
		if err := newTimeSession.SetState(tsession.StateWork); err != nil {
			return err
		}
		if err := newTimeSession.CreatedUpdated.SetCreatedBy(actor.BasicUser); err != nil {
			return err
		}
		if err := newTimeSession.CreatedUpdated.SetUpdatedBy(actor.BasicUser); err != nil {
			return err
		}
		e.openTimeSession = newTimeSession
	}

	if err := e.openTimeSession.StartWorking(inc); err != nil {
		return err
	}

	return nil
}

func (e *FieldEngineer) canStartWorking(actor actor.Actor) error {
	if !actor.IsFieldEngineer() {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "actor is not field engineer")
	}

	if actor.IsFieldEngineer() && *actor.FieldEngineerID() != e.uuid {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "actor is not this field engineer")
	}

	return nil
}

func (e *FieldEngineer) canStartTravelling(actor actor.Actor) error {
	if !actor.IsFieldEngineer() {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "actor is not field engineer")
	}

	if actor.IsFieldEngineer() && *actor.FieldEngineerID() != e.uuid {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "actor is not this field engineer")
	}

	if e.HasOpenTimeSession() {
		return domain.NewErrorf(domain.ErrorCodeActionForbidden, "actor already has an open time session")
	}

	return nil
}
