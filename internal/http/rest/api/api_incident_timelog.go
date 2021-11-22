package api

// Timelog object
// swagger:model
type Timelog struct {
	// required: true
	UUID UUID `json:"uuid"`

	// required: true
	Remote bool `json:"remote"`

	// Time spent working in seconds
	// minimum: 0
	Work uint `json:"work,omitempty"`

	VisitSummary string `json:"visit_summary,omitempty"`

	CreatedUpdated
}
