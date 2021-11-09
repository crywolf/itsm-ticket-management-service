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

// CreatedBy ...
func (o CreatedUpdated) CreatedBy() ref.ExternalUserUUID {
	return o.createdBy
}

// CreatedAt ...
func (o CreatedUpdated) CreatedAt() DateTime {
	return o.createdAt
}

// UpdatedBy ...
func (o CreatedUpdated) UpdatedBy() ref.ExternalUserUUID {
	return o.updatedBy
}

// UpdatedAt ...
func (o CreatedUpdated) UpdatedAt() DateTime {
	return o.updatedAt
}

// SetCreatedBy returns error if createdBy was already set
func (o *CreatedUpdated) SetCreatedBy(userID ref.ExternalUserUUID) error {
	if !o.createdBy.IsZero() {
		return errors.New("cannot set CreatedBy, it was already set")
	}
	o.createdBy = userID
	return nil
}

// SetCreatedAt returns error if createdAt was already set
func (o *CreatedUpdated) SetCreatedAt(dateTime DateTime) error {
	if !o.createdAt.IsZero() {
		return errors.New("cannot set CreatedAt, it was already set")
	}

	o.createdAt = dateTime
	return nil
}

// SetUpdatedBy ...
func (o *CreatedUpdated) SetUpdatedBy(userID ref.ExternalUserUUID) error {
	o.updatedBy = userID
	return nil
}

// SetUpdatedAt ...
func (o *CreatedUpdated) SetUpdatedAt(dateTime DateTime) error {
	o.updatedAt = dateTime
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
