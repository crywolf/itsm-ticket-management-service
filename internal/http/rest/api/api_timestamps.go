package api

// CreatedUpdated contains timestamps and user who created/updated the resource
type CreatedUpdated struct {
	CreatedInfo
	UpdatedInfo
}

// CreatedInfo contains timestamp and user who created the resource
type CreatedInfo struct {
	// Time when the resource was created
	// required: true
	// swagger:strfmt date-time
	CreatedAt string `json:"created_at"`

	// Reference to the user who created this resource
	// required: true
	// swagger:strfmt uuid
	CreatedBy string `json:"created_by"`
}

// UpdatedInfo contains timestamp and user who updated the resource
type UpdatedInfo struct {
	// Time when the resource was updated
	// swagger:strfmt date-time
	UpdatedAt string `json:"updated_at"`

	// Reference to the user who updated this resource
	// swagger:strfmt uuid
	UpdatedBy string `json:"updated_by"`
}
