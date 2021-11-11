package rest

import "github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"

func (s *Server) presenters() {
	s.presenter = presenters.NewIncidentPresenter(s.logger, s.ExternalLocationAddress)
}
