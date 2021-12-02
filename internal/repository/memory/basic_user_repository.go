package memory

import (
	"context"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
)

// BasicUserRepositoryMemory keeps data in memory
type BasicUserRepositoryMemory struct {
	//Rand      io.Reader
	//Clock     Clock
	users []user.BasicUser
}

// GetBasicUser returns the Basic User with given ID from the repository
func (r *BasicUserRepositoryMemory) GetBasicUser(_ context.Context, _ ref.ChannelID, ID ref.UUID) (user.BasicUser, error) {
	return user.BasicUser{}, nil
}

// GetBasicUserByExternalID returns the Basic User with given external ID from the repository
func (r *BasicUserRepositoryMemory) GetBasicUserByExternalID(_ context.Context, _ ref.ChannelID, externalID ref.ExternalUserUUID) (user.BasicUser, error) {
	return user.BasicUser{}, nil
}
