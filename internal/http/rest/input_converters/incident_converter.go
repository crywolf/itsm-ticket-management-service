package converters

import (
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/input_converters/validators"
	"go.uber.org/zap"
)

// NewIncidentPayloadConverter creates an incident input payload converting service
func NewIncidentPayloadConverter(logger *zap.SugaredLogger, validator validators.PayloadValidator) IncidentPayloadConverter {
	return &incidentPayloadConverter{
		BasePayloadConverter: NewBasePayloadConverter(logger, validator),
	}
}

type incidentPayloadConverter struct {
	*BasePayloadConverter
}

// IncidentCreateParamsFromBody converts JSON payload to api.CreateIncidentParams
func (c incidentPayloadConverter) IncidentCreateParamsFromBody(r *http.Request) (api.CreateIncidentParams, error) {
	var payload api.CreateIncidentParams

	if err := c.unmarshalFromBody(r, &payload); err != nil {
		return payload, err
	}

	return payload, nil
}

// IncidentUpdateParamsFromBody converts JSON payload to api.UpdateIncidentParams
func (c incidentPayloadConverter) IncidentUpdateParamsFromBody(r *http.Request) (api.UpdateIncidentParams, error) {
	var payload api.UpdateIncidentParams

	if err := c.unmarshalFromBody(r, &payload); err != nil {
		return payload, err
	}

	return payload, nil
}

// IncidentStartWorkingParamsFromBody converts JSON payload to api.IncidentStartWorkingParams
func (c incidentPayloadConverter) IncidentStartWorkingParamsFromBody(r *http.Request) (api.IncidentStartWorkingParams, error) {
	var payload api.IncidentStartWorkingParams

	if err := c.unmarshalFromBody(r, &payload); err != nil {
		return payload, err
	}

	return payload, nil
}
