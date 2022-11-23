package rest

import (
	converters "github.com/crywolf/itsm-ticket-management-service/internal/http/rest/input_converters"
	"github.com/crywolf/itsm-ticket-management-service/internal/http/rest/input_converters/validators"
)

type jsonInputPayloadConverters struct {
	incident converters.IncidentPayloadConverter
}

func (s *Server) registerInputConverters() {
	validator := validators.NewPayloadValidator()

	s.inputPayloadConverters.incident = converters.NewIncidentPayloadConverter(s.logger, validator)
}
