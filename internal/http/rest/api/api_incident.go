package api

import "github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"

// Incident API object
// swagger:model
type Incident struct {
	// required: true
	// swagger:strfmt uuid
	UUID string `json:"uuid"`

	// Unique identifier provided by user creating the incident
	// required: true
	Number string `json:"number"`

	// ID in external system
	ExternalID string `json:"external_id,omitempty"`

	// required: true
	ShortDescription string `json:"short_description"`

	Description string `json:"description,omitempty"`

	// State of the ticket
	// required: true
	// example: new
	State incident.State `json:"state"`

	// List of timelogs
	// read only: true
	Timelogs []UUID `json:"timelogs,omitempty"`

	CreatedUpdated
}

// CreateIncidentParams is the payload used to create new incident
// swagger:model
type CreateIncidentParams struct {
	// Unique identifier
	// required: true
	Number string `json:"number" validate:"required"`

	// ID in external system
	ExternalID string `json:"external_id"`

	// required: true
	ShortDescription string `json:"short_description" validate:"required"`

	Description string `json:"description"`
}

// IncidentResponse ...
type IncidentResponse struct {
	Incident
	Links HypermediaLinks `json:"_links,omitempty"`
}

// Data structure representing a single incident
// swagger:response incidentResponse
type incidentResponseWrapper struct {
	// in: body
	Body struct {
		IncidentResponse
	}
}

// IncidentListResponse ...
type IncidentListResponse struct {
	// Total number of elements in the list
	// required: true
	Total int `json:"total"`
	// Number of elements on the current page
	// required: true
	Size int `json:"size"`
	// Current page number
	// required: true
	Page   int                `json:"page"`
	Result []IncidentResponse `json:"_embedded,omitempty"`
	// example: {self:{href:example.com}}
	Links HypermediaLinks `json:"_links,omitempty"`
}

// A list of incidents
// swagger:response incidentListResponse
type incidentListResponseWrapper struct {
	// in: body
	Body struct {
		IncidentListResponse
	}
}

// Created
// swagger:response incidentCreatedResponse
type incidentCreatedResponseWrapper struct {
	// URI of the resource
	// example: http://localhost:8080/incidents/2af4f493-0bd5-4513-b440-6cbb465feadb
	// in: header
	Location string
}

// swagger:parameters CreateIncident
type createIncidentParameterWrapper struct {
	AuthorizationHeaders

	// in: body
	// required: true
	Body CreateIncidentParams
}
