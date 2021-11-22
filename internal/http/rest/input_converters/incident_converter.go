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

// IncidentParamsFromBody converts JSON payload to api.CreateIncidentParams
func (c incidentPayloadConverter) IncidentParamsFromBody(r *http.Request) (api.CreateIncidentParams, error) {
	var incPayload api.CreateIncidentParams

	if err := c.unmarshalFromBody(r, &incPayload); err != nil {
		return incPayload, err
	}

	return incPayload, nil
}
