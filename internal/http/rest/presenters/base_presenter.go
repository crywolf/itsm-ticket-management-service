package presenters

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters/hypermedia"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
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

// RenderLocationHeader sends Location header containing URI in the form 'route/resourceID'.
// Use it for rendering location of newly created resource
func (p BasePresenter) RenderLocationHeader(w http.ResponseWriter, route string, resourceID ref.UUID) {
	resourceURI := fmt.Sprintf("%s%s/%s", p.serverAddr, route, resourceID)

	w.Header().Set("Location", resourceURI)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// RenderError replies to the request with the specified error message and HTTP code
func (p BasePresenter) RenderError(w http.ResponseWriter, msg string, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		if msg == "" {
			msg = httpErr.Error()
		}
		p.renderErrorJSON(w, msg, httpErr.Code())
		return
	}

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

	p.renderErrorJSON(w, msg, status)
}

func (p BasePresenter) resourceToHypermediaLinks(domainObject hypermedia.ActionsMapper, hypermediaMapper hypermedia.Mapper, inList bool) api.HypermediaLinks {
	hypermediaLinks := api.HypermediaLinks{}

	actions := hypermediaMapper.RoutesToHypermediaActionLinks()
	allowedActions := domainObject.AllowedActions(hypermediaMapper.Actor())
	for _, action := range allowedActions {
		link := actions[action]
		href := strings.ReplaceAll(link.Href, "{uuid}", domainObject.UUID().String())
		hypermediaLinks[link.Name] = map[string]string{
			"href": href,
		}
	}

	if inList {
		hypermediaLinks.AppendSelfLink(fmt.Sprintf("%s%s/%s", p.serverAddr, hypermediaMapper.RequestURL().Path, domainObject.UUID()))
	} else {
		hypermediaLinks.AppendSelfLink(hypermediaMapper.SelfLink())
	}

	return hypermediaLinks
}

func (p BasePresenter) hypermediaListLinks(hypermediaMapper hypermedia.Mapper, pagination *repository.Pagination) api.HypermediaListLinks {
	hypermediaLinks := api.HypermediaListLinks{}

	query := hypermediaMapper.RequestURL().Query()
	url := *hypermediaMapper.RequestURL()
	first := url
	last := url
	prev := url
	next := url

	if pagination.First == 1 {
		query.Del("page")
	} else {
		query.Set("page", strconv.Itoa(pagination.First))
	}
	first.RawQuery = query.Encode()

	if pagination.Last == 1 {
		query.Del("page")
	} else {
		query.Set("page", strconv.Itoa(pagination.Last))
	}
	last.RawQuery = query.Encode()

	prevString := ""
	if pagination.Prev == 1 {
		query.Del("page")
		prev.RawQuery = query.Encode()
		prevString = fmt.Sprintf("%s%s", p.serverAddr, prev.String())
	} else if pagination.Prev > 1 {
		query.Set("page", strconv.Itoa(pagination.Prev))
		prev.RawQuery = query.Encode()
		prevString = fmt.Sprintf("%s%s", p.serverAddr, prev.String())
	}

	nextString := ""
	if pagination.Next > 0 {
		query.Set("page", strconv.Itoa(pagination.Next))
		next.RawQuery = query.Encode()
		nextString = fmt.Sprintf("%s%s", p.serverAddr, next.String())
	}

	hypermediaLinks["first"] = map[string]string{
		"href": fmt.Sprintf("%s%s", p.serverAddr, first.String()),
	}
	hypermediaLinks["last"] = map[string]string{
		"href": fmt.Sprintf("%s%s", p.serverAddr, last.String()),
	}

	if prevString != "" {
		hypermediaLinks["prev"] = map[string]string{
			"href": prevString,
		}
	}
	if nextString != "" {
		hypermediaLinks["next"] = map[string]string{
			"href": nextString,
		}
	}

	hypermediaLinks.AppendSelfLink(hypermediaMapper.SelfLink())

	return hypermediaLinks
}

// renderErrorJSON replies to the request with the specified error message and HTTP code.
// It encodes error string as JSON object {"error":"error_string"} and sets correct header.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
// The error message should be plain text.
func (p BasePresenter) renderErrorJSON(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	errorJSON, err := json.Marshal(msg)
	if err != nil {
		p.logger.Errorw("sending json error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(code)
	_, _ = fmt.Fprintf(w, `{"error":%s}`+"\n", errorJSON)
}

// renderJSON encodes 'v' to JSON and writes it to the 'w'. Also sets correct Content-Type header.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
func (p BasePresenter) renderJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		err = domain.WrapErrorf(err, domain.ErrorCodeUnknown, "could not encode JSON response")
		p.logger.Errorw("encoding json", "error", err)
		p.renderErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
