package mocks

import (
	"context"

	fieldengineer "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/stretchr/testify/mock"
)

// FieldEngineerServiceMock is a field engineer service mock
type FieldEngineerServiceMock struct {
	mock.Mock
}

// GetFieldEngineer mock
func (s *FieldEngineerServiceMock) GetFieldEngineer(_ context.Context, channelID ref.ChannelID, actor actor.Actor, ID ref.UUID) (fieldengineer.FieldEngineer, error) {
	args := s.Called(channelID, actor, ID)
	return args.Get(0).(fieldengineer.FieldEngineer), args.Error(1)
}
