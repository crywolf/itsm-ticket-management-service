package fieldengineersvc

import (
	"context"

	fieldengineer "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
)

// FieldEngineerService provides incident operations
type FieldEngineerService interface {

	// GetFieldEngineer returns the field engineer with the given ID from the repository
	GetFieldEngineer(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, ID ref.UUID) (fieldengineer.FieldEngineer, error)
}
