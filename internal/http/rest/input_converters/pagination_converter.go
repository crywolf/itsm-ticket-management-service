package converters

import (
	"net/http"
	"strconv"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
)

type paginationParams struct {
	page         uint
	itemsPerPage uint
}

// NewPaginationParams parses request query and returns params with information about requested page and items per page to be displayed
func NewPaginationParams(r *http.Request, actorUser actor.Actor) (PaginationParams, error) {
	var err error
	var page64 uint64

	queryValues := r.URL.Query()
	pageParam := queryValues.Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	if page64, err = strconv.ParseUint(pageParam, 10, 0); err != nil {
		return nil, presenters.NewErrorf(http.StatusBadRequest, "incorrect 'page' parameter: '%s'", pageParam)
	}
	if page64 == 0 {
		return nil, presenters.NewErrorf(http.StatusBadRequest, "incorrect 'page' parameter: '%s'", pageParam)
	}
	return &paginationParams{
		page:         uint(page64),
		itemsPerPage: actorUser.BasicUser.ItemsPerPage(),
	}, nil
}

func (p paginationParams) Page() uint {
	return p.page
}

func (p paginationParams) ItemsPerPage() uint {
	return p.itemsPerPage
}
