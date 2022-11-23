package api

// Timelog object
// swagger:model
type Timelog struct {
	// required: true
	UUID UUID `json:"uuid"`

	// required: true
	Remote bool `json:"remote"`

	// Time when the timelog was created
	// swagger:strfmt date-time
	Start string `json:"start,omitempty"`

	// Time when the timelog was closed
	// swagger:strfmt date-time
	End string `json:"end,omitempty"`

	// Time spent working in seconds
	// minimum: 0
	Work uint `json:"work,omitempty"`

	VisitSummary string `json:"visit_summary,omitempty"`

	CreatedUpdated
}

// TimelogResponse ...
type TimelogResponse struct {
	Timelog
	Links    HypermediaLinks   `json:"_links,omitempty"`
	Embedded EmbeddedResources `json:"_embedded,omitempty"`
}

// Data structure representing a single timelog
// swagger:response timelogResponse
type timelogResponseWrapper struct {
	// in: body
	Body struct {
		TimelogResponse
	}
}

// swagger:parameters GetIncidentTimelog
type generalTimelogParameterWrapper struct {
	AuthorizationHeaders

	// ID of the resource
	// in: path
	// required: true
	TicketUUID UUID `json:"uuid"`

	// ID of the timelog
	// in: path
	// required: true
	UUID UUID `json:"timelog_uuid"`
}
