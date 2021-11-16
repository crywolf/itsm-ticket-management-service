package rest

import (
	converters "github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/input_converters"
)

type jsonInputPayloadConverters struct {
	incident  converters.IncidentPayloadConverter
}

func (s *Server) registerInputConverters() {
	s.inputPayloadConverters.incident = converters.NewIncidentPayloadConverter(s.logger)
}
