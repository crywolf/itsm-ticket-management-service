package converters

import (
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
)

// IncidentPayloadConverter provides conversion from JSON request body payload to object
type IncidentPayloadConverter interface {
	// IncidentParamsFromBody converts JSON payload to api.CreateIncidentParams
	IncidentParamsFromBody(r *http.Request) (api.CreateIncidentParams, error)
}
