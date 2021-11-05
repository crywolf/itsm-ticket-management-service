package incident

import (
	"errors"

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

	createdUpdated types.CreatedUpdated
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

func (i Incident) GetCreatedBy() ref.ExternalUserUUID {
	return i.createdUpdated.CreatedBy
}

func (i *Incident) SetCreatedBy(userID ref.ExternalUserUUID) error {
	if !i.createdUpdated.CreatedBy.IsZero() {
		return errors.New("cannot set CreatedBy, it was already set")
	}
	i.createdUpdated.CreatedBy = userID
	return nil
}

func (i Incident) GetUpdatedBy() ref.ExternalUserUUID {
	return i.createdUpdated.UpdatedBy
}

func (i *Incident) SetUpdatedBy(userID ref.ExternalUserUUID) error {
	i.createdUpdated.UpdatedBy = userID
	return nil
}

func (i Incident) GetCreatedAt() types.DateTime {
	return i.createdUpdated.CreatedAt
}

func (i *Incident) SetCreatedAt(dateTime types.DateTime) error {
	if !i.createdUpdated.CreatedAt.IsZero() {
		return errors.New("cannot set CreatedAt, it was already set")
	}

	i.createdUpdated.CreatedAt = dateTime
	return nil
}

func (i Incident) GetUpdatedAt() types.DateTime {
	return i.createdUpdated.CreatedAt
}

func (i *Incident) SetUpdatedAt(dateTime types.DateTime) error {
	i.createdUpdated.UpdatedAt = dateTime
	return nil
}

func (i Incident) GetCreatedUpdated() types.CreatedUpdated {
	return i.createdUpdated
}
