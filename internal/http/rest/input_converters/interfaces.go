package converters

import (
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
)

// PaginationParams provides information about current requested page number and a number of items per page to be displayed
type PaginationParams interface {
	// Page is the requested page number to be returned
	Page() uint

	// ItemsPerPage returns how many items per page should be displayed
	ItemsPerPage() uint
}

// IncidentPayloadConverter provides conversion from JSON request body payload to object
type IncidentPayloadConverter interface {
	// IncidentCreateParamsFromBody converts JSON payload to api.CreateIncidentParams
	IncidentCreateParamsFromBody(r *http.Request) (api.CreateIncidentParams, error)

	// IncidentUpdateParamsFromBody converts JSON payload to api.UpdateIncidentParams
	IncidentUpdateParamsFromBody(r *http.Request) (api.UpdateIncidentParams, error)

	// IncidentStartWorkingParamsFromBody converts JSON payload to api.IncidentStartWorkingParams
	IncidentStartWorkingParamsFromBody(r *http.Request) (api.IncidentStartWorkingParams, error)

	// IncidentStopWorkingParamsFromBody converts JSON payload to api.IncidentStopWorkingParams
	IncidentStopWorkingParamsFromBody(r *http.Request) (api.IncidentStopWorkingParams, error)
}
