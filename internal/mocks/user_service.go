package mocks

import (
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/stretchr/testify/mock"
)

// UserServiceMock is a user service mock
type UserServiceMock struct {
	mock.Mock
}

// ActorFromRequest returns Actor object that represents a user who initiated the request
func (s *UserServiceMock) ActorFromRequest(r *http.Request) (actor.Actor, error) {
	args := s.Called(r)
	return args.Get(0).(actor.Actor), args.Error(1)
}
