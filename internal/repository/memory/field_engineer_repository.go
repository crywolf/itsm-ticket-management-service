package memory

import (
	"context"
	"io"

	"github.com/crywolf/itsm-ticket-management-service/internal/domain"
	fieldengineer "github.com/crywolf/itsm-ticket-management-service/internal/domain/field_engineer"
	tsession "github.com/crywolf/itsm-ticket-management-service/internal/domain/field_engineer/time_session"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/types"
	"github.com/crywolf/itsm-ticket-management-service/internal/repository"
)

// FieldEngineerRepositoryMemory keeps data in memory
type FieldEngineerRepositoryMemory struct {
	basicUserRepository repository.BasicUserRepository
	Rand                io.Reader
	clock               repository.Clock
	fieldEngineers      []FieldEngineer
	timeSessions        map[string]TimeSession
}

// NewFieldEngineerRepositoryMemory returns new initialized repository
func NewFieldEngineerRepositoryMemory(clock repository.Clock, basicUserRepo repository.BasicUserRepository) repository.FieldEngineerRepository {
	return &FieldEngineerRepositoryMemory{
		basicUserRepository: basicUserRepo,
		clock:               clock,
		timeSessions:        make(map[string]TimeSession),
	}
}

// AddFieldEngineer adds the given field engineer to the repository
func (r *FieldEngineerRepositoryMemory) AddFieldEngineer(_ context.Context, _ ref.ChannelID, fe fieldengineer.FieldEngineer) (ref.UUID, error) {
	now := r.clock.NowFormatted().String()

	feID, err := repository.GenerateUUID(r.Rand)
	if err != nil {
		return ref.UUID(""), err
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

// UpdateFieldEngineer updates the given field engineer in the repository
func (r *FieldEngineerRepositoryMemory) UpdateFieldEngineer(_ context.Context, _ ref.ChannelID, fe fieldengineer.FieldEngineer) (ref.UUID, error) {
	var err error
	now := r.clock.NowFormatted().String()

	if fe.HasOpenTimeSession() {
		openTS := fe.OpenTimeSession()
		createdAt := openTS.CreatedUpdated.CreatedAt().String()
		updatedAt := openTS.CreatedUpdated.UpdatedAt().String()

		tSessionID := openTS.UUID()
		if tSessionID.IsZero() { // newly opened => set new UUID
			tSessionID, err = repository.GenerateUUID(r.Rand)
			if err != nil {
				return ref.UUID(""), err
			}

			createdAt = now
			updatedAt = now

			fe.TimeSessions = append(fe.TimeSessions, tSessionID)
		}

		var incidents []IncidentInfo
		for _, incInfo := range openTS.Incidents {
			incidents = append(incidents, IncidentInfo{
				ID:                 incInfo.IncidentID.String(),
				HasSupplierProduct: incInfo.HasSupplierProduct,
			})
		}

		storedTS := TimeSession{
			ID:                          tSessionID.String(),
			State:                       openTS.State().String(),
			Incidents:                   incidents,
			Work:                        openTS.Work,
			Travel:                      openTS.Travel,
			TravelBack:                  openTS.TravelBack,
			TravelDistanceInTravelUnits: openTS.TravelDistanceInTravelUnits,
			CreatedAt:                   createdAt,
			CreatedBy:                   openTS.CreatedUpdated.CreatedByID().String(),
			UpdatedAt:                   updatedAt,
			UpdatedBy:                   openTS.CreatedUpdated.UpdatedByID().String(),
		}

		r.timeSessions[storedTS.ID] = storedTS
	}

	var tsUUIDs []string
	for _, tsID := range fe.TimeSessions {
		tsUUIDs = append(tsUUIDs, tsID.String())
	}

	storedFE := FieldEngineer{
		ID:           fe.UUID().String(),
		BasicUserID:  fe.BasicUser.UUID().String(),
		TimeSessions: tsUUIDs,
		CreatedBy:    fe.CreatedUpdated.CreatedByID().String(),
		CreatedAt:    fe.CreatedUpdated.CreatedAt().String(),
		UpdatedBy:    fe.CreatedUpdated.UpdatedByID().String(),
		UpdatedAt:    now,
	}

	for i := range r.fieldEngineers {
		if r.fieldEngineers[i].ID == fe.UUID().String() {
			r.fieldEngineers[i] = storedFE
			return fe.UUID(), nil
		}
	}

	return fe.UUID(), domain.WrapErrorf(ErrNotFound, domain.ErrorCodeNotFound, "error updating field engineer in repository")
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

	// set Time sessions (UUIDs)
	var tsUUIDs []ref.UUID
	for _, tsID := range storedFE.TimeSessions {
		tsUUIDs = append(tsUUIDs, ref.UUID(tsID))
	}

	fe.TimeSessions = tsUUIDs

	// load and set open time session if any
	openTS, err := r.loadOpenTimeSession(ctx, channelID, storedFE)
	if err != nil {
		return fieldengineer.FieldEngineer{}, err
	}
	fe.SetOpenTimeSession(openTS)

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

// loadOpenTimeSession loads field engineer's open time session if any
func (r FieldEngineerRepositoryMemory) loadOpenTimeSession(ctx context.Context, channelID ref.ChannelID, storedFE FieldEngineer) (*tsession.TimeSession, error) {
	errMsg := "error loading field engineer from repository (%s)"

	for _, tsID := range storedFE.TimeSessions {
		storedTS := r.timeSessions[tsID]

		var incidents []tsession.IncidentInfo
		for _, incInfo := range storedTS.Incidents {
			incidents = append(incidents, tsession.IncidentInfo{
				IncidentID:         ref.UUID(incInfo.ID),
				HasSupplierProduct: incInfo.HasSupplierProduct,
			})
		}

		if storedTS.State != "closed" { // time session is open
			openTS := &tsession.TimeSession{
				Incidents:                   incidents,
				Work:                        storedTS.Work,
				Travel:                      storedTS.Travel,
				TravelBack:                  storedTS.TravelBack,
				TravelDistanceInTravelUnits: storedTS.TravelDistanceInTravelUnits,
			}

			state, err := tsession.NewStateFromString(storedTS.State)
			if err != nil {
				return nil, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedTimeSession.state")
			}

			if err := openTS.SetState(state); err != nil {
				return nil, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedTimeSession.SetState")
			}

			createdByUser, err := r.basicUserRepository.GetBasicUser(ctx, channelID, ref.UUID(storedTS.CreatedBy))
			if err != nil {
				return nil, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedTimeSession.createdBy")
			}

			err = openTS.CreatedUpdated.SetCreated(createdByUser, types.DateTime(storedTS.CreatedAt))
			if err != nil {
				return nil, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedTimeSession.createdAt")
			}

			updatedByUser, err := r.basicUserRepository.GetBasicUser(ctx, channelID, ref.UUID(storedTS.UpdatedBy))
			if err != nil {
				return nil, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedTimeSession.UpdatedBy")
			}

			err = openTS.CreatedUpdated.SetUpdated(updatedByUser, types.DateTime(storedTS.UpdatedAt))
			if err != nil {
				return nil, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedTimeSession.UpdatedAt")
			}

			return openTS, nil
		}
	}

	return nil, nil
}
