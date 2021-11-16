package converters

import (
	"encoding/json"
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
	"go.uber.org/zap"
)

// NewIncidentPayloadConverter creates an incident input payload converting service
func NewIncidentPayloadConverter(logger *zap.SugaredLogger) IncidentPayloadConverter {
	return &incidentPayloadConverter{
		BasePayloadConverter: NewBasePayloadConverter(logger),
	}
}

type incidentPayloadConverter struct {
	*BasePayloadConverter
}

// IncidentParamsFromBody converts JSON payload to api.CreateIncidentParams
func (c incidentPayloadConverter) IncidentParamsFromBody(r *http.Request) (api.CreateIncidentParams, error) {
	var incPayload api.CreateIncidentParams

	reqBody, err := c.readBody(r)
	if err != nil {
		return incPayload, err
	}

	err = json.Unmarshal(reqBody, &incPayload)
	if err != nil {
		err = presenters.WrapErrorf(err, http.StatusBadRequest, "could not decode JSON from request")
		return incPayload, err
	}

	return incPayload, nil
}
