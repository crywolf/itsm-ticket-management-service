package memory

import (
	"context"
	"testing"

	fieldengineer "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident/timelog"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
	"github.com/KompiTech/itsm-ticket-management-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncidentRepositoryMemory_AddingAndGettingIncident(t *testing.T) {
	basicUser := user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	err := basicUser.SetUUID("f49d5fd5-8da4-4779-b5ba-32e78aa2c444")
	require.NoError(t, err)

	clock := mocks.NewFixedClock()
	basicUserRepository := &BasicUserRepositoryMemory{
		users: []user.BasicUser{basicUser},
	}
	fieldEngineerRepository := NewFieldEngineerRepositoryMemory(clock, basicUserRepository)

	repo := NewIncidentRepositoryMemory(clock, basicUserRepository, fieldEngineerRepository)

	channelID := ref.ChannelID("e27ddcd0-0e1f-4bc5-93df-f6f04155beec")
	ctx := context.Background()

	inc1 := incident.Incident{
		Number:           "ABC123",
		ExternalID:       "some external ID",
		ShortDescription: "some short description",
		Description:      "some description",
	}
	err = inc1.CreatedUpdated.SetCreatedBy(basicUser)
	require.NoError(t, err)
	err = inc1.CreatedUpdated.SetUpdatedBy(basicUser)
	require.NoError(t, err)

	incID, err := repo.AddIncident(ctx, channelID, inc1)
	require.NoError(t, err)

	retInc, err := repo.GetIncident(ctx, channelID, incID)
	require.NoError(t, err)

	assert.Equal(t, incID, retInc.UUID())
	assert.Equal(t, inc1.Number, retInc.Number)
	assert.Equal(t, inc1.ExternalID, retInc.ExternalID)
	assert.Equal(t, inc1.ShortDescription, retInc.ShortDescription)
	assert.Equal(t, inc1.Description, retInc.Description)
	assert.Empty(t, retInc.FieldEngineerID)
	assert.Empty(t, retInc.Timelogs)

	// test correct timestamps
	assert.NotEmpty(t, inc1.CreatedUpdated.CreatedByID())
	assert.Equal(t, inc1.CreatedUpdated.CreatedByID(), retInc.CreatedUpdated.CreatedByID())
	assert.Equal(t, clock.NowFormatted(), retInc.CreatedUpdated.CreatedAt())

	assert.NotEmpty(t, inc1.CreatedUpdated.UpdatedBy())
	assert.Equal(t, inc1.CreatedUpdated.UpdatedBy(), retInc.CreatedUpdated.UpdatedBy())
	assert.Equal(t, clock.NowFormatted(), retInc.CreatedUpdated.UpdatedAt())
}

func TestIncidentRepositoryMemory_UpdateIncident(t *testing.T) {
	basicUser := user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	err := basicUser.SetUUID("f49d5fd5-8da4-4779-b5ba-32e78aa2c444")
	require.NoError(t, err)

	basicUser2 := user.BasicUser{
		ExternalUserUUID: "ee824cad-d7a6-4f48-87dc-e8461a9201c4",
		Name:             "Jan",
		Surname:          "Novak",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	err = basicUser2.SetUUID("00271cb4-3716-4203-9124-1d2f515ae0b2")
	require.NoError(t, err)

	basicUserRepository := &BasicUserRepositoryMemory{
		users: []user.BasicUser{basicUser, basicUser2},
	}

	fieldEngineerUUID := ref.UUID("1adb8393-cff0-489c-a82f-3fe5d15708d4")
	fieldEngineer := fieldengineer.FieldEngineer{
		BasicUser: basicUser,
	}
	err = fieldEngineer.SetUUID(fieldEngineerUUID)
	require.NoError(t, err)

	err = fieldEngineer.CreatedUpdated.SetCreatedBy(basicUser)
	require.NoError(t, err)

	err = fieldEngineer.CreatedUpdated.SetUpdatedBy(basicUser)
	require.NoError(t, err)

	storedFieldEngineer := FieldEngineer{
		ID:          fieldEngineerUUID.String(),
		BasicUserID: basicUser.UUID().String(),
		CreatedBy:   basicUser.UUID().String(),
		UpdatedBy:   basicUser.UUID().String(),
	}

	fieldEngineerRepository := &FieldEngineerRepositoryMemory{
		basicUserRepository: basicUserRepository,
		fieldEngineers:      []FieldEngineer{storedFieldEngineer},
	}

	clock := mocks.NewFixedClock()
	repo := NewIncidentRepositoryMemory(clock, basicUserRepository, fieldEngineerRepository)

	channelID := ref.ChannelID("e27ddcd0-0e1f-4bc5-93df-f6f04155beec")
	ctx := context.Background()

	inc1 := incident.Incident{
		Number:           "ABC123",
		ExternalID:       "some external ID",
		ShortDescription: "some short description",
		Description:      "some description",
	}
	err = inc1.CreatedUpdated.SetCreatedBy(basicUser)
	require.NoError(t, err)
	err = inc1.CreatedUpdated.SetUpdatedBy(basicUser)
	require.NoError(t, err)

	incID, err := repo.AddIncident(ctx, channelID, inc1)
	require.NoError(t, err)

	retInc, err := repo.GetIncident(ctx, channelID, incID)
	require.NoError(t, err)

	changedDescription := "some changed description"
	retInc.Description = changedDescription
	retInc.FieldEngineerID = &fieldEngineerUUID
	err = retInc.SetState(incident.StateInProgress)
	require.NoError(t, err)

	// set open timelog
	openTimelog := &timelog.Timelog{}
	err = openTimelog.CreatedUpdated.SetCreated(basicUser2, clock.NowFormatted())
	require.NoError(t, err)
	err = openTimelog.CreatedUpdated.SetUpdated(basicUser, clock.NowFormatted())
	require.NoError(t, err)

	retInc.SetOpenTimelog(openTimelog)

	// update incident with open timelog
	retIncID, err := repo.UpdateIncident(ctx, channelID, retInc)
	require.NoError(t, err)
	assert.Equal(t, incID, retIncID)

	// get updated incident
	updatedInc, err := repo.GetIncident(ctx, channelID, incID)
	require.NoError(t, err)

	assert.Equal(t, incID, updatedInc.UUID())
	assert.Equal(t, inc1.Number, updatedInc.Number)
	assert.Equal(t, inc1.ExternalID, updatedInc.ExternalID)
	assert.Equal(t, inc1.ShortDescription, updatedInc.ShortDescription)
	assert.Equal(t, changedDescription, updatedInc.Description)
	assert.Equal(t, fieldEngineer.UUID(), *updatedInc.FieldEngineerID)
	assert.Len(t, updatedInc.Timelogs, 1, "timelogs count")

	assert.Equal(t, retInc.OpenTimelog(), updatedInc.OpenTimelog())
}

func TestIncidentRepositoryMemory_ListIncidents(t *testing.T) {
	basicUser := user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	_ = basicUser.SetUUID("f49d5fd5-8da4-4779-b5ba-32e78aa2c444")

	basicUser2 := user.BasicUser{
		ExternalUserUUID: "ee824cad-d7a6-4f48-87dc-e8461a9201c4",
		Name:             "Jan",
		Surname:          "Novak",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	err := basicUser2.SetUUID("00271cb4-3716-4203-9124-1d2f515ae0b2")
	require.NoError(t, err)

	clock := mocks.NewFixedClock()
	basicUserRepository := &BasicUserRepositoryMemory{
		users: []user.BasicUser{basicUser, basicUser2},
	}
	fieldEngineerRepository := NewFieldEngineerRepositoryMemory(clock, basicUserRepository)

	repo := NewIncidentRepositoryMemory(clock, basicUserRepository, fieldEngineerRepository)

	channelID := ref.ChannelID("e27ddcd0-0e1f-4bc5-93df-f6f04155beec")

	ctx := context.Background()

	// empty list
	emptyList, err := repo.ListIncidents(ctx, channelID, 1, 10)
	require.NoError(t, err)

	// pagination
	assert.Equal(t, 0, emptyList.Size)
	assert.Equal(t, 0, emptyList.Total)
	assert.Equal(t, 1, emptyList.Page)
	assert.Equal(t, 1, emptyList.First)
	assert.Equal(t, 1, emptyList.Last)
	assert.Equal(t, 0, emptyList.Prev)
	assert.Equal(t, 0, emptyList.Next)
	assert.Len(t, emptyList.Result, 0)

	inc1 := incident.Incident{
		Number:           "Bca258",
		ExternalID:       "some external ID",
		ShortDescription: "some short description",
		Description:      "some description",
	}
	err = inc1.CreatedUpdated.SetCreatedBy(basicUser)
	require.NoError(t, err)
	err = inc1.CreatedUpdated.SetUpdatedBy(basicUser)
	require.NoError(t, err)

	inc2 := incident.Incident{
		Number:           "CDB36478",
		ExternalID:       "some external ID 2",
		ShortDescription: "some short description 2",
		Description:      "some description 2",
	}
	err = inc2.CreatedUpdated.SetCreatedBy(basicUser)
	require.NoError(t, err)
	err = inc2.CreatedUpdated.SetUpdatedBy(basicUser2)
	require.NoError(t, err)

	_, err = repo.AddIncident(ctx, channelID, inc1)
	require.NoError(t, err)

	_, err = repo.AddIncident(ctx, channelID, inc2)
	require.NoError(t, err)

	// first page
	incidentsList, err := repo.ListIncidents(ctx, channelID, 1, 10)
	require.NoError(t, err)

	// pagination
	assert.Equal(t, 2, incidentsList.Size)
	assert.Equal(t, 2, incidentsList.Total)
	assert.Equal(t, 1, incidentsList.Page)
	assert.Equal(t, 1, incidentsList.First)
	assert.Equal(t, 1, incidentsList.Last)
	assert.Equal(t, 0, incidentsList.Prev)
	assert.Equal(t, 0, incidentsList.Next)

	list := incidentsList.Result

	assert.Len(t, list, 2)

	for i, retInc := range list {
		var inc incident.Incident
		switch i {
		case 0:
			inc = inc1
		case 1:
			inc = inc2
		}

		assert.NotEmpty(t, retInc.UUID())
		assert.Equal(t, inc.Number, retInc.Number)
		assert.Equal(t, inc.ExternalID, retInc.ExternalID)
		assert.Equal(t, inc.ShortDescription, retInc.ShortDescription)
		assert.Equal(t, inc.Description, retInc.Description)

		// test correct timestamps
		assert.NotEmpty(t, inc.CreatedUpdated.CreatedBy())
		assert.Equal(t, inc.CreatedUpdated.CreatedBy(), retInc.CreatedUpdated.CreatedBy())
		assert.Equal(t, clock.NowFormatted(), retInc.CreatedUpdated.CreatedAt())

		assert.NotEmpty(t, inc.CreatedUpdated.UpdatedBy())
		assert.Equal(t, inc.CreatedUpdated.UpdatedBy(), retInc.CreatedUpdated.UpdatedBy())
		assert.Equal(t, clock.NowFormatted(), retInc.CreatedUpdated.UpdatedAt())
	}

	// second page out of range
	incidentsList, err = repo.ListIncidents(ctx, channelID, 2, 10)
	require.NoError(t, err)

	list = incidentsList.Result
	assert.Len(t, list, 0)

	// pagination
	assert.Equal(t, 0, incidentsList.Size)
	assert.Equal(t, 2, incidentsList.Total)
	assert.Equal(t, 2, incidentsList.Page)
	assert.Equal(t, 1, incidentsList.First)
	assert.Equal(t, 1, incidentsList.Last)
	assert.Equal(t, 1, incidentsList.Prev)
	assert.Equal(t, 0, incidentsList.Next)

	// first page with small number per page
	incidentsList, err = repo.ListIncidents(ctx, channelID, 1, 1)
	require.NoError(t, err)

	// pagination
	assert.Equal(t, 1, incidentsList.Size)
	assert.Equal(t, 2, incidentsList.Total)
	assert.Equal(t, 1, incidentsList.Page)
	assert.Equal(t, 1, incidentsList.First)
	assert.Equal(t, 2, incidentsList.Last)
	assert.Equal(t, 0, incidentsList.Prev)
	assert.Equal(t, 2, incidentsList.Next)

	list = incidentsList.Result
	assert.Len(t, list, 1)

	// second page with small number per page
	incidentsList, err = repo.ListIncidents(ctx, channelID, 2, 1)
	require.NoError(t, err)

	// pagination
	assert.Equal(t, 1, incidentsList.Size)
	assert.Equal(t, 2, incidentsList.Total)
	assert.Equal(t, 2, incidentsList.Page)
	assert.Equal(t, 1, incidentsList.First)
	assert.Equal(t, 2, incidentsList.Last)
	assert.Equal(t, 1, incidentsList.Prev)
	assert.Equal(t, 0, incidentsList.Next)

	list = incidentsList.Result
	assert.Len(t, list, 1)
}
