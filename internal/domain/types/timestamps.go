package types

import "github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"

// CreatedUpdated contains timestamps and user who created/updated the resource
type CreatedUpdated struct {
	createdInfo
	updatedInfo
}

// createdInfo contains timestamp and user who created the resource
type createdInfo struct {
	// Time when the resource was created
	CreatedAt DateTime

	// Reference to the user who created this resource
	CreatedBy ref.ExternalUserUUID
}

// updatedInfo contains timestamp and user who updated the resource
type updatedInfo struct {
	// Time when the resource was updated
	UpdatedAt DateTime

	// Reference to the user who updated this resource
	UpdatedBy ref.ExternalUserUUID
}
