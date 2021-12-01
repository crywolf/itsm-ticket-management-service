package rest

import (
	"context"
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	usersvc "github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/service"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
	grpc2http "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc/status"
)

type userKeyType int

var userKey userKeyType

// AddUserInfo is a middleware that stores info about invoking user in request context
// (or about user this request is made on behalf of)
func (s Server) AddUserInfo(next httprouter.Handle, us usersvc.Service) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		actorUser, err := us.ActorFromRequest(r)
		if err != nil {
			s.logger.Error("AddUserInfo middleware: ActorFromRequest failed:", "error", err)
			httpStatusCode := grpc2http.HTTPStatusFromCode(status.Code(err))
			err := presenters.WrapErrorf(err, httpStatusCode, "could not retrieve correct user info from user service")
			s.presenters.base.RenderError(w, "", err)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, &actorUser)

		next(w, r.WithContext(ctx), ps)
	}
}

// ActorFromContext returns the Actor user stored in request's context if any.
func (s Server) ActorFromContext(ctx context.Context) (actor.Actor, bool) {
	u, ok := ctx.Value(userKey).(*actor.Actor)
	return *u, ok
}
