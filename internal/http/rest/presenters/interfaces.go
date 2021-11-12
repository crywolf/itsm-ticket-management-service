package presenters

import (
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters/hypermedia"
)

// ErrorPresenter allows replying with error
type ErrorPresenter interface {
	// WriteError replies to the request with the specified error message and HTTP code.
	// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
	// The error message should be plain text.
	WriteError(w http.ResponseWriter, error string, code int)
}

// IncidentPresenter provides REST responses for incident resource
type IncidentPresenter interface {
	ErrorPresenter

	// WriteIncident encodes incident and writes it to 'w'.  Also sets correct Content-Type header.
	// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
	WriteIncident(w http.ResponseWriter, incident incident.Incident, hypermediaMapper hypermedia.Mapper)

	// WriteIncidentList encodes list of incidents and writes it to 'w'.  Also sets correct Content-Type header.
	// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
	WriteIncidentList(w http.ResponseWriter, incidentList []incident.Incident, hypermediaMapper hypermedia.Mapper)
}
