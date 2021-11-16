package presenters

import (
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters/hypermedia"
)

// ErrorPresenter allows replying with error
type ErrorPresenter interface {
	// RenderError replies to the request with the specified error message and HTTP code.
	// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
	// The error message should be plain text.
	RenderError(w http.ResponseWriter, msg string, err error)
}

// IncidentPresenter provides REST responses for incident resource
type IncidentPresenter interface {
	ErrorPresenter

	// RenderIncident encodes incident and writes it to 'w'.  Also sets correct Content-Type header.
	// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
	RenderIncident(w http.ResponseWriter, incident incident.Incident, hypermediaMapper hypermedia.Mapper)

	// RenderIncidentList encodes list of incidents and writes it to 'w'.  Also sets correct Content-Type header.
	// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
	RenderIncidentList(w http.ResponseWriter, incidentList []incident.Incident, hypermediaMapper hypermedia.Mapper)
}
