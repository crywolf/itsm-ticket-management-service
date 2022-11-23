package mocks

import (
	"context"

	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/stretchr/testify/mock"
)

// ExternalUserServiceMock is a user service mock
type ExternalUserServiceMock struct {
	mock.Mock
}

// ActorFromRequest returns an Actor object that represents a user who initiated the request
func (s *ExternalUserServiceMock) ActorFromRequest(ctx context.Context, authToken string, channelID ref.ChannelID, onBehalf string) (actor.Actor, error) {
	args := s.Called(authToken, channelID, onBehalf)
	return args.Get(0).(actor.Actor), args.Error(1)
}
