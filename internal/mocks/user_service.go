package mocks

import (
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/stretchr/testify/mock"
)

// UserServiceMock is a user service mock
type UserServiceMock struct {
	mock.Mock
}

// ActorFromRequest returns an Actor object that represents a user who initiated the request
func (s *UserServiceMock) ActorFromRequest(authToken string, channelID ref.ChannelID, onBehalf string) (actor.Actor, error) {
	args := s.Called(authToken, channelID, onBehalf)
	return args.Get(0).(actor.Actor), args.Error(1)
}
