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

// swagger:parameters CreateIncident ListIncidents
type generalNoParameterWrapper struct {
	AuthorizationHeaders
}

// swagger:parameters GetIncident UpdateIncident IncidentStartWorking
type generalIDParameterWrapper struct {
	AuthorizationHeaders

	// ID of the resource
	// in: path
	// required: true
	UUID UUID `json:"uuid"`
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

// ActionLink represents action link to be transformed to HAL hypermedia links
type ActionLink struct {
	Name string
	Href string
}

// PageInfo represents information about data set in the list
type PageInfo struct {
	// Total number of elements in the list
	// required: true
	Total int `json:"total"`
	// Size of dataset of elements on the current page
	// required: true
	Size int `json:"size"`
	// Current page number
	// required: true
	Page int `json:"page"`
}

// HypermediaLinks contain links to other API calls
// example: {"self": {"href": "example.com"}}
type HypermediaLinks map[string]interface{}

// AppendSelfLink adds resource's 'self' link
func (l *HypermediaLinks) AppendSelfLink(url string) {
	(*l)["self"] = map[string]string{
		"href": url,
	}
}

// Link represents HAL hypermedia link
type Link struct {
	// swagger:strfmt uri
	Href string `json:"href"`
}

// HypermediaListLinks contain 'self' and pagination links to be use in list views
type HypermediaListLinks struct {
	// Self link
	// required: true
	Self Link `json:"self"`
	// First page link
	// required: true
	First Link `json:"first"`
	// Last page link
	// required: true
	Last Link `json:"last"`
	// Previous page link
	Prev *Link `json:"prev,omitempty"`
	// Next page link
	Next *Link `json:"next,omitempty"`
}

// AppendSelfLink adds resource's 'self' link
func (l *HypermediaListLinks) AppendSelfLink(url string) {
	l.Self = Link{Href: url}
}

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
