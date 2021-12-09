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

// CreatedBy getter // TODO just for testing - dodělat
func (o CreatedUpdated) CreatedBy() *user.BasicUser {
	basicUser := &user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Kolecko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	_ = basicUser.SetUUID("cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0")

	return basicUser
}

// CreatedByID getter
func (o CreatedUpdated) CreatedByID() ref.UUID {
	return o.createdBy
}

// CreatedAt getter
func (o CreatedUpdated) CreatedAt() DateTime {
	return o.createdAt
}

// UpdatedBy getter
func (o CreatedUpdated) UpdatedBy() ref.UUID {
	return o.updatedBy
}

// UpdatedAt getter
func (o CreatedUpdated) UpdatedAt() DateTime {
	return o.updatedAt
}

// SetCreated sets info about the user and time when the resource was created. It  returns error if createdBy or createdAt was already set.
func (o *CreatedUpdated) SetCreated(userID ref.UUID, dateTime DateTime) error {
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
func (o *CreatedUpdated) SetUpdated(userID ref.UUID, dateTime DateTime) error {
	if err := o.SetUpdatedBy(userID); err != nil {
		return err
	}
	o.updatedAt = dateTime
	return nil
}

// SetCreatedBy sets info about the user who crested this resource. It returns error if createdBy was already set.
func (o *CreatedUpdated) SetCreatedBy(userID ref.UUID) error {
	if !o.createdBy.IsZero() {
		return errors.New("cannot set CreatedByID, it was already set")
	}
	o.createdBy = userID
	return nil
}

// SetUpdatedBy sets info about the user who updated this resource. It returns error if createdBy was already set.
func (o *CreatedUpdated) SetUpdatedBy(userID ref.UUID) error {
	o.updatedBy = userID
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

	// Reference to the user who created this resource
	createdBy ref.UUID
}

// updatedInfo contains timestamp and user who updated the resource
type updatedInfo struct {
	// Time when the resource was updated
	updatedAt DateTime

	// Reference to the user who updated this resource
	updatedBy ref.UUID
}
