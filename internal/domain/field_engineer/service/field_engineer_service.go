package fieldengineersvc

import (
	"context"

	fieldengineer "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
)

// NewFieldEngineerService creates the field engineer service
func NewFieldEngineerService(repo repository.FieldEngineerRepository) FieldEngineerService {
	return &fieldEngineerService{repo}
}

type fieldEngineerService struct {
	repo repository.FieldEngineerRepository
}

func (s *fieldEngineerService) GetFieldEngineer(ctx context.Context, channelID ref.ChannelID, actor actor.Actor, ID ref.UUID) (fieldengineer.FieldEngineer, error) {
	return s.repo.GetFieldEngineer(ctx, channelID, ID)
}
