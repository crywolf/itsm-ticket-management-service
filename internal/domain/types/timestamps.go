package types

import (
	"errors"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/embedded"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
)

// CreatedUpdated contains timestamps and reference to user who created/updated the resource
type CreatedUpdated struct {
	createdInfo
	updatedInfo
}

// CreatedBy returns the user who created the resource
func (o CreatedUpdated) CreatedBy() user.BasicUser {
	return o.createdBy
}

// CreatedByID returns ID of the user who created the resource
func (o CreatedUpdated) CreatedByID() ref.UUID {
	return o.createdBy.UUID()
}

// CreatedAt getter
func (o CreatedUpdated) CreatedAt() DateTime {
	return o.createdAt
}

// UpdatedBy returns the user who updated the resource
func (o CreatedUpdated) UpdatedBy() user.BasicUser {
	return o.updatedBy
}

// UpdatedByID returns ID of the user who updated the resource
func (o CreatedUpdated) UpdatedByID() ref.UUID {
	return o.updatedBy.UUID()
}

// UpdatedAt getter
func (o CreatedUpdated) UpdatedAt() DateTime {
	return o.updatedAt
}

// SetCreated sets info about the user and time when the resource was created. It  returns error if createdBy or createdAt was already set.
func (o *CreatedUpdated) SetCreated(basicUser user.BasicUser, dateTime DateTime) error {
	if err := o.SetCreatedBy(basicUser); err != nil {
		return err
	}
	if !o.createdAt.IsZero() {
		return errors.New("cannot set CreatedAt, it was already set")
	}
	o.createdAt = dateTime
	return nil
}

// SetUpdated sets info about the user and time when the resource was updated
func (o *CreatedUpdated) SetUpdated(basicUser user.BasicUser, dateTime DateTime) error {
	if err := o.SetUpdatedBy(basicUser); err != nil {
		return err
	}
	o.updatedAt = dateTime
	return nil
}

// SetCreatedBy sets info about the user who crested this resource. It returns error if createdBy was already set.
func (o *CreatedUpdated) SetCreatedBy(basicUser user.BasicUser) error {
	if !o.createdBy.IsZero() {
		return errors.New("cannot set CreatedByID, it was already set")
	}
	o.createdBy = basicUser
	return nil
}

// SetUpdatedBy sets info about the user who updated this resource. It returns error if createdBy was already set.
func (o *CreatedUpdated) SetUpdatedBy(basicUser user.BasicUser) error {
	o.updatedBy = basicUser
	return nil
}

// EmbeddedResources should be called from the same func in the resource that includes CreatedUpdated object
func (o CreatedUpdated) EmbeddedResources(_ actor.Actor) []embedded.Resource {
	resources := []embedded.Resource{
		embedded.CreatedBy,
		embedded.UpdatedBy,
	}
	return resources
}

// createdInfo contains timestamp and user who created the resource
type createdInfo struct {
	// Time when the resource was created
	createdAt DateTime

	// User who created this resource
	createdBy user.BasicUser
}

// updatedInfo contains timestamp and user who updated the resource
type updatedInfo struct {
	// Time when the resource was updated
	updatedAt DateTime

	// User who updated this resource
	updatedBy user.BasicUser
}
