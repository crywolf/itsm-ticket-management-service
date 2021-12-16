package memory

import (
	"context"
	"io"
	"log"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	fieldengineer "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
)

// FieldEngineerRepositoryMemory keeps data in memory
type FieldEngineerRepositoryMemory struct {
	basicUserRepository repository.BasicUserRepository
	Rand                io.Reader
	clock               Clock
	fieldEngineers      []FieldEngineer
	timeSessions        map[string]Timelog
}

// NewFieldEngineerRepositoryMemory returns new initialized repository
func NewFieldEngineerRepositoryMemory(clock Clock, basicUserRepo repository.BasicUserRepository) repository.FieldEngineerRepository {
	return &FieldEngineerRepositoryMemory{
		basicUserRepository: basicUserRepo,
		clock:               clock,
		timeSessions:        make(map[string]Timelog),
	}
}

// AddFieldEngineer adds the given field engineer to the repository
func (r *FieldEngineerRepositoryMemory) AddFieldEngineer(_ context.Context, _ ref.ChannelID, fe fieldengineer.FieldEngineer) (ref.UUID, error) {
	now := r.clock.NowFormatted().String()

	feID, err := repository.GenerateUUID(r.Rand)
	if err != nil {
		log.Fatal(err)
	}

	storedFE := FieldEngineer{
		ID:          feID.String(),
		BasicUserID: fe.BasicUser.UUID().String(),
		CreatedBy:   fe.CreatedUpdated.CreatedByID().String(),
		CreatedAt:   now,
		UpdatedBy:   fe.CreatedUpdated.UpdatedByID().String(),
		UpdatedAt:   now,
	}

	r.fieldEngineers = append(r.fieldEngineers, storedFE)

	return feID, nil
}

// GetFieldEngineer returns the field engineer with given ID from the repository
func (r *FieldEngineerRepositoryMemory) GetFieldEngineer(ctx context.Context, channelID ref.ChannelID, ID ref.UUID) (fieldengineer.FieldEngineer, error) {
	var inc fieldengineer.FieldEngineer
	var err error

	for i := range r.fieldEngineers {
		if r.fieldEngineers[i].ID == ID.String() {
			storedFE := r.fieldEngineers[i]

			inc, err = r.convertStoredToDomainFieldEngineer(ctx, channelID, storedFE)
			if err != nil {
				return fieldengineer.FieldEngineer{}, err
			}

			return inc, nil
		}
	}

	return fieldengineer.FieldEngineer{}, domain.WrapErrorf(ErrNotFound, domain.ErrorCodeNotFound, "error loading field engineer from repository")
}

func (r FieldEngineerRepositoryMemory) convertStoredToDomainFieldEngineer(ctx context.Context, channelID ref.ChannelID, storedFE FieldEngineer) (fieldengineer.FieldEngineer, error) {
	var fe fieldengineer.FieldEngineer
	errMsg := "error loading field engineer from repository (%s)"

	err := fe.SetUUID(ref.UUID(storedFE.ID))
	if err != nil {
		return fieldengineer.FieldEngineer{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedFE.ID")
	}

	basicUser, err := r.basicUserRepository.GetBasicUser(ctx, channelID, ref.UUID(storedFE.BasicUserID))
	if err != nil {
		return fieldengineer.FieldEngineer{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedFE.BasicUser")
	}

	fe.BasicUser = basicUser

	// TODO load and set open time session if any

	createdByUser, err := r.basicUserRepository.GetBasicUser(ctx, channelID, ref.UUID(storedFE.CreatedBy))
	if err != nil {
		return fieldengineer.FieldEngineer{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedFE.CreatedBy")
	}

	err = fe.CreatedUpdated.SetCreated(createdByUser, types.DateTime(storedFE.CreatedAt))
	if err != nil {
		return fieldengineer.FieldEngineer{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedFE.CreatedAt")
	}

	updatedByUser, err := r.basicUserRepository.GetBasicUser(ctx, channelID, ref.UUID(storedFE.UpdatedBy))
	if err != nil {
		return fieldengineer.FieldEngineer{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedFE.UpdatedBy")
	}

	err = fe.CreatedUpdated.SetUpdated(updatedByUser, types.DateTime(storedFE.UpdatedAt))
	if err != nil {
		return fieldengineer.FieldEngineer{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedFE.UpdatedAt")
	}

	return fe, nil
}
