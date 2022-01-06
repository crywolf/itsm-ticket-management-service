package memory

import (
	"context"
	"testing"

	fieldengineer "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer"
	tsession "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer/time_session"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
	"github.com/KompiTech/itsm-ticket-management-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldEngineerRepositoryMemory_AddingAndGettingFieldEngineer(t *testing.T) {
	adminBasicUser := user.BasicUser{
		ExternalUserUUID: "2d839741-da07-4256-bd53-4030bb0effeb",
		Name:             "Admin",
		Surname:          "User",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	err := adminBasicUser.SetUUID("b8d49f19-5e54-44cf-b547-f16bacb69294")
	require.NoError(t, err)

	feBasicUser := user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	err = feBasicUser.SetUUID("f49d5fd5-8da4-4779-b5ba-32e78aa2c444")
	require.NoError(t, err)

	clock := mocks.NewFixedClock()
	basicUserRepository := &BasicUserRepositoryMemory{
		users: []user.BasicUser{adminBasicUser, feBasicUser},
	}
	repo := NewFieldEngineerRepositoryMemory(clock, basicUserRepository)

	channelID := ref.ChannelID("e27ddcd0-0e1f-4bc5-93df-f6f04155beec")
	ctx := context.Background()

	fe := fieldengineer.FieldEngineer{BasicUser: feBasicUser}
	err = fe.CreatedUpdated.SetCreatedBy(adminBasicUser)
	require.NoError(t, err)
	err = fe.CreatedUpdated.SetUpdatedBy(adminBasicUser)
	require.NoError(t, err)

	feID, err := repo.AddFieldEngineer(ctx, channelID, fe)
	require.NoError(t, err)

	retFe, err := repo.GetFieldEngineer(ctx, channelID, feID)
	require.NoError(t, err)

	assert.Equal(t, feID, retFe.UUID())
	assert.Equal(t, fe.BasicUser, retFe.BasicUser)
	assert.Len(t, fe.TimeSessions, 0, "time sessions count")

	// test correct timestamps
	assert.NotEmpty(t, fe.CreatedUpdated.CreatedByID())
	assert.Equal(t, fe.CreatedUpdated.CreatedBy(), retFe.CreatedUpdated.CreatedBy())
	assert.Equal(t, clock.NowFormatted(), retFe.CreatedUpdated.CreatedAt())

	assert.NotEmpty(t, fe.CreatedUpdated.UpdatedByID())
	assert.Equal(t, fe.CreatedUpdated.UpdatedBy(), retFe.CreatedUpdated.UpdatedBy())
	assert.Equal(t, clock.NowFormatted(), retFe.CreatedUpdated.UpdatedAt())
}

func TestFieldEngineerRepositoryMemory_UpdateFieldEngineer(t *testing.T) {
	adminBasicUser := user.BasicUser{
		ExternalUserUUID: "2d839741-da07-4256-bd53-4030bb0effeb",
		Name:             "Admin",
		Surname:          "User",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	err := adminBasicUser.SetUUID("b8d49f19-5e54-44cf-b547-f16bacb69294")
	require.NoError(t, err)

	feBasicUser := user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	err = feBasicUser.SetUUID("f49d5fd5-8da4-4779-b5ba-32e78aa2c444")
	require.NoError(t, err)

	clock := mocks.NewFixedClock()
	basicUserRepository := &BasicUserRepositoryMemory{
		users: []user.BasicUser{adminBasicUser, feBasicUser},
	}
	repo := NewFieldEngineerRepositoryMemory(clock, basicUserRepository)

	channelID := ref.ChannelID("e27ddcd0-0e1f-4bc5-93df-f6f04155beec")
	ctx := context.Background()

	fe := fieldengineer.FieldEngineer{BasicUser: feBasicUser}
	err = fe.CreatedUpdated.SetCreatedBy(adminBasicUser)
	require.NoError(t, err)
	err = fe.CreatedUpdated.SetUpdatedBy(adminBasicUser)
	require.NoError(t, err)

	feID, err := repo.AddFieldEngineer(ctx, channelID, fe)
	require.NoError(t, err)

	retFe, err := repo.GetFieldEngineer(ctx, channelID, feID)
	require.NoError(t, err)

	// set open time session
	openTS := tsession.TimeSession{
		Incidents: []tsession.IncidentInfo{{
			IncidentID:         "d67c7799-cab5-4dbd-8a5c-2e4e19070f77",
			HasSupplierProduct: false,
		}},
		Work: 3600,
	}
	err = openTS.SetState(tsession.StateWork)
	require.NoError(t, err)
	err = openTS.CreatedUpdated.SetCreated(feBasicUser, clock.NowFormatted())
	require.NoError(t, err)
	err = openTS.CreatedUpdated.SetUpdated(feBasicUser, clock.NowFormatted())
	require.NoError(t, err)

	retFe.SetOpenTimeSession(&openTS)

	// update field engineer with open time session
	retFeID, err := repo.UpdateFieldEngineer(ctx, channelID, retFe)
	require.NoError(t, err)
	assert.Equal(t, feID, retFeID)

	// get updated field engineer
	updatedFe, err := repo.GetFieldEngineer(ctx, channelID, feID)
	require.NoError(t, err)

	assert.Len(t, updatedFe.TimeSessions, 1, "time sessions count")
	assert.Equal(t, retFe.OpenTimeSession(), updatedFe.OpenTimeSession())
	assert.Len(t, retFe.OpenTimeSession().Incidents, 1)
	assert.Equal(t, retFe.OpenTimeSession().Incidents, updatedFe.OpenTimeSession().Incidents)
}
