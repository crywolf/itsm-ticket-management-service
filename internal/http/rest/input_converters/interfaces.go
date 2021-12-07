package converters

import (
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
)

// PaginationParams provides information about current requested page number and a number of items per page to be displayed
type PaginationParams interface {
	// Page returns requested page number to be returned
	Page() uint

	// ItemsPerPage returns how many items per page should be displayed
	ItemsPerPage() uint
}

// IncidentPayloadConverter provides conversion from JSON request body payload to object
type IncidentPayloadConverter interface {
	// IncidentParamsFromBody converts JSON payload to api.CreateIncidentParams
	IncidentParamsFromBody(r *http.Request) (api.CreateIncidentParams, error)
}
