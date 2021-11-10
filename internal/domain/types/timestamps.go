package types

import (
	"errors"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
)

// CreatedUpdated contains timestamps and reference to user who created/updated the resource
type CreatedUpdated struct {
	createdInfo
	updatedInfo
}

// CreatedBy getter
func (o CreatedUpdated) CreatedBy() ref.ExternalUserUUID {
	return o.createdBy
}

// CreatedAt getter
func (o CreatedUpdated) CreatedAt() DateTime {
	return o.createdAt
}

// UpdatedBy getter
func (o CreatedUpdated) UpdatedBy() ref.ExternalUserUUID {
	return o.updatedBy
}

// UpdatedAt getter
func (o CreatedUpdated) UpdatedAt() DateTime {
	return o.updatedAt
}

// SetCreated sets info about the user and time when the resource was created. It  returns error if createdBy or createdAt was already set
func (o *CreatedUpdated) SetCreated(userID ref.ExternalUserUUID, dateTime DateTime) error {
	if err := o.SetCreatedBy(userID); err != nil {
		return err
	}
	if !o.createdAt.IsZero() {
		return errors.New("cannot set CreatedAt, it was already set")
	}
	o.createdAt = dateTime
	return nil
}

// SetUpdated sets info about the user and time when the resource was updated
func (o *CreatedUpdated) SetUpdated(userID ref.ExternalUserUUID, dateTime DateTime) error {
	if err := o.SetUpdatedBy(userID); err != nil {
		return err
	}
	o.updatedAt = dateTime
	return nil
}

// SetCreatedBy sets info about the user who crested this resource. It returns error if createdBy was already set
func (o *CreatedUpdated) SetCreatedBy(userID ref.ExternalUserUUID) error {
	if !o.createdBy.IsZero() {
		return errors.New("cannot set CreatedBy, it was already set")
	}
	o.createdBy = userID
	return nil
}

// SetUpdatedBy sets info about the user who updated this resource. It returns error if createdBy was already set
func (o *CreatedUpdated) SetUpdatedBy(userID ref.ExternalUserUUID) error {
	o.updatedBy = userID
	return nil
}

// createdInfo contains timestamp and user who created the resource
type createdInfo struct {
	// Time when the resource was created
	createdAt DateTime

	// Reference to the user who created this resource
	createdBy ref.ExternalUserUUID
}

// updatedInfo contains timestamp and user who updated the resource
type updatedInfo struct {
	// Time when the resource was updated
	updatedAt DateTime

	// Reference to the user who updated this resource
	updatedBy ref.ExternalUserUUID
}
