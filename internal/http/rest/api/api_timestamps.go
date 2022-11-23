package api

import (
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/types"
)

// CreatedUpdated contains timestamps and user who created/updated the resource
type CreatedUpdated struct {
	CreatedInfo
	UpdatedInfo
}

// NewCreatedUpdatedInfo creates new initialized CreatedUpdatedInfo
func NewCreatedUpdatedInfo(cu types.CreatedUpdated) CreatedUpdated {
	return CreatedUpdated{
		CreatedInfo: NewCreatedInfo(cu),
		UpdatedInfo: NewUpdatedInfo(cu),
	}
}

// CreatedInfo contains timestamp and user who created the resource
type CreatedInfo struct {
	// Time when the resource was created
	// required: true
	// swagger:strfmt date-time
	CreatedAt string `json:"created_at,omitempty"`

	// Reference to the user who created this resource
	// required: true
	// swagger:strfmt uuid
	CreatedBy string `json:"created_by,omitempty"`
}

// NewCreatedInfo creates new initialized CreatedInfo
func NewCreatedInfo(cu types.CreatedUpdated) CreatedInfo {
	return CreatedInfo{
		CreatedAt: cu.CreatedAt().String(),
		CreatedBy: cu.CreatedByID().String(),
	}
}

// UpdatedInfo contains timestamp and user who updated the resource
type UpdatedInfo struct {
	// Time when the resource was updated
	// swagger:strfmt date-time
	UpdatedAt string `json:"updated_at,omitempty"`

	// Reference to the user who updated this resource
	// swagger:strfmt uuid
	UpdatedBy string `json:"updated_by,omitempty"`
}

// NewUpdatedInfo creates new initialized UpdatedInfo
func NewUpdatedInfo(cu types.CreatedUpdated) UpdatedInfo {
	return UpdatedInfo{
		UpdatedAt: cu.UpdatedAt().String(),
		UpdatedBy: cu.UpdatedByID().String(),
	}
}
