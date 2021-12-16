package api

// FieldEngineer API object
// swagger:model
type FieldEngineer struct {
	// required: true
	// swagger:strfmt uuid
	UUID string `json:"uuid"`

	// required: true
	BasicUser BasicUser `json:"basic_user"`
}
