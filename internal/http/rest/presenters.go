package rest

import (
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
)

type jsonPresenters struct {
	base *presenters.BasePresenter
	incident  presenters.IncidentPresenter
}

func (s *Server) registerPresenters() {
	s.presenters.base = presenters.NewBasePresenter(s.logger, s.ExternalLocationAddress)
	s.presenters.incident = presenters.NewIncidentPresenter(s.logger, s.ExternalLocationAddress)
}
