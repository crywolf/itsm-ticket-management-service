package rest

import (
	converters "github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/input_converters"
)

func (s *Server) registerInputConverters() {
	s.inputPayloadConverter = converters.NewIncidentPayloadConverter(s.logger)
}
