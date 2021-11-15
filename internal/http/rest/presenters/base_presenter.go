package presenters

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters/hypermedia"
	"go.uber.org/zap"
)

// NewBasePresenter returns presentation service with basic functionality
func NewBasePresenter(logger *zap.SugaredLogger, serverAddr string) *BasePresenter {
	return &BasePresenter{
		logger:     logger,
		serverAddr: serverAddr,
	}
}

// BasePresenter must be included in all derived presenters via object composition
type BasePresenter struct {
	logger     *zap.SugaredLogger
	serverAddr string
}

// WriteError replies to the request with the specified error message and HTTP code
func (p BasePresenter) WriteError(w http.ResponseWriter, msg string, err error) {
	status := http.StatusInternalServerError

	var dErr *domain.Error
	if !errors.As(err, &dErr) {
		msg = fmt.Sprintf("internal error: %s", err.Error())
	} else {
		if msg == "" {
			msg = dErr.Error()
		}

		switch dErr.Code() {
		case domain.ErrorCodeInvalidArgument:
			status = http.StatusBadRequest
		case domain.ErrorCodeNotFound:
			status = http.StatusNotFound
		case domain.ErrorCodeUserNotAuthorized:
			status = http.StatusUnauthorized
		case domain.ErrorCodeActionForbidden:
			status = http.StatusForbidden
		case domain.ErrorCodeUnknown:
			fallthrough
		default:
			status = http.StatusInternalServerError
		}
	}

	p.sendErrorJSON(w, msg, status)
}

func (p BasePresenter) resourceToHypermediaLinks(hypermediaMapper hypermedia.Mapper, domainObject hypermedia.ActionsMapper) api.HypermediaLinks {
	hypermediaLinks := api.HypermediaLinks{}

	actions := hypermediaMapper.RoutesToHypermediaActionLinks()
	allowedActions := domainObject.AllowedActions()
	for _, action := range allowedActions {
		link := actions[action]
		href := strings.ReplaceAll(link.Href, "{uuid}", domainObject.UUID().String())
		hypermediaLinks[link.Name] = map[string]string{
			"href": href,
		}
	}

	return hypermediaLinks
}

// sendErrorJSON replies to the request with the specified error message and HTTP code.
// It encodes error string as JSON object {"error":"error_string"} and sets correct header.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
// The error message should be plain text.
func (p BasePresenter) sendErrorJSON(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	errorJSON, err := json.Marshal(error)
	if err != nil {
		p.logger.Errorw("sending json error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(code)
	_, _ = fmt.Fprintf(w, `{"error":%s}`+"\n", errorJSON)
}

// encodeJSON encodes 'v' to JSON and writes it to the 'w'. Also sets correct Content-Type header.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
func (p BasePresenter) encodeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		err = domain.WrapErrorf(err, domain.ErrorCodeUnknown, "could not encode JSON response")
		p.logger.Errorw("encoding json", "error", err)
		p.sendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
