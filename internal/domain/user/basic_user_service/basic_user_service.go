package basicusersvc

import (
	"context"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
)

// BasicUserService provides Basic User operations
type BasicUserService interface {
	// GetBasicUser returns the Basic User with given ID from the repository
	GetBasicUser(ctx context.Context, channelID ref.ChannelID, ID ref.UUID) (user.BasicUser, error)

	// GetBasicUserByExternalID returns the Basic User with given external ID from the repository
	GetBasicUserByExternalID(ctx context.Context, channelID ref.ChannelID, externalID ref.ExternalUserUUID) (user.BasicUser, error)
}

// BasicUserRepository provides access to the Basic User repository
type BasicUserRepository interface {
	// GetBasicUser returns the Basic User with given ID from the repository
	GetBasicUser(ctx context.Context, channelID ref.ChannelID, ID ref.UUID) (user.BasicUser, error)

	// GetBasicUserByExternalID returns the Basic User with given external ID from the repository
	GetBasicUserByExternalID(ctx context.Context, channelID ref.ChannelID, externalID ref.ExternalUserUUID) (user.BasicUser, error)
}

// NewBasicUserService creates the Basic User service
func NewBasicUserService(r BasicUserRepository) BasicUserService {
	return &service{r}
}

func (s *service) GetBasicUser(ctx context.Context, channelID ref.ChannelID, ID ref.UUID) (user.BasicUser, error) {
	return s.r.GetBasicUser(ctx, channelID, ID)
}

func (s *service) GetBasicUserByExternalID(ctx context.Context, channelID ref.ChannelID, externalID ref.ExternalUserUUID) (user.BasicUser, error) {
	return s.r.GetBasicUserByExternalID(ctx, channelID, externalID)
}

type service struct {
	r BasicUserRepository
}
