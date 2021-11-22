// Package api ITSM Ticket Management Service REST API
//
// Documentation for ITSM Ticket Management Service REST API
//
//	Schemes: http
//	BasePath: /
//	Version: 0.0.1
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package api

// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handlers

// UUID represents UUID of a resource
// swagger:strfmt uuid
type UUID string

// swagger:parameters GetIncident
type generalIDParameterWrapper struct {
	AuthorizationHeaders

	// ID of the resource
	// in: path
	// required: true
	UUID UUID `json:"uuid"`
}

// swagger:parameters ListIncidents
type generalListParameterWrapper struct {
	AuthorizationHeaders
}

// AuthorizationHeaders represents general authorization header parameters used in many API calls
type AuthorizationHeaders struct {
	// Bearer token
	// in: header
	// required: true
	Authorization string `json:"authorization"`

	// in: header
	// required: true
	ChannelID UUID `json:"channel-id"`
}

// ActionLink represents HAL hypermedia links
type ActionLink struct {
	Name string
	Href string
}

// HypermediaLinks contain links to other API calls
type HypermediaLinks map[string]interface{}

// No content
// swagger:response deleteNoContentResponse
type deleteNoContentResponseWrapper struct{}

// Error
// swagger:response errorResponse
type errorResponseWrapper struct {
	// in: body
	Body struct {
		// required: true
		// Description of the error
		ErrorMessage string `json:"error"`
	}
}

// Bad Request
// swagger:response errorResponse400
type errorResponseWrapper400 errorResponseWrapper

// Unauthorized
// swagger:response errorResponse401
type errorResponseWrapper401 errorResponseWrapper

// Forbidden
// swagger:response errorResponse403
type errorResponseWrapper403 errorResponseWrapper

// Not Found
// swagger:response errorResponse404
type errorResponseWrapper404 errorResponseWrapper

// Conflict
// swagger:response errorResponse409
type errorResponseWrapper409 errorResponseWrapper
