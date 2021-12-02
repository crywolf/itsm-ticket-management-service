package memory

import (
	"context"
	"io"
	"log"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
)

// BasicUserRepositoryMemory keeps data in memory
type BasicUserRepositoryMemory struct {
	Rand  io.Reader
	users []user.BasicUser
}

// AddBasicUser adds the Basic User to the repository
func (r *BasicUserRepositoryMemory) AddBasicUser(_ context.Context, _ ref.ChannelID, user user.BasicUser) (ref.UUID, error) {
	id, err := repository.GenerateUUID(r.Rand)
	if err != nil {
		log.Fatal(err)
	}

	r.users = append(r.users, user)

	return id, nil

}

// GetBasicUser returns the Basic User with the given ID from the repository
func (r *BasicUserRepositoryMemory) GetBasicUser(_ context.Context, _ ref.ChannelID, ID ref.UUID) (user.BasicUser, error) {
	for i := range r.users {
		if r.users[i].UUID().String() == ID.String() {
			storedUser := r.users[i]
			return storedUser, nil
		}
	}
	return user.BasicUser{}, domain.WrapErrorf(ErrNotFound, domain.ErrorCodeNotFound, "repo GetBasicUser")
}

// GetBasicUserByExternalID returns the Basic User with the given external ID from the repository
func (r *BasicUserRepositoryMemory) GetBasicUserByExternalID(_ context.Context, _ ref.ChannelID, externalID ref.ExternalUserUUID) (user.BasicUser, error) {
	for i := range r.users {
		if r.users[i].ExternalUserUUID == externalID {
			storedUser := r.users[i]
			return storedUser, nil
		}
	}
	return user.BasicUser{}, domain.WrapErrorf(ErrNotFound, domain.ErrorCodeNotFound, "repo GetBasicUserByExternalID")
}
