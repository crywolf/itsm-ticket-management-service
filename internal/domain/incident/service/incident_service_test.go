package incidentsvc

import (
	"context"
	"testing"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/KompiTech/itsm-ticket-management-service/internal/mocks"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_incidentService_CreateAndRetrieveIncidents(t *testing.T) {
	channelID := ref.ChannelID("e27ddcd0-0e1f-4bc5-93df-f6f04155beec")
	ctx := context.Background()

	basicUser := user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}

	basicUserRepository := &memory.BasicUserRepositoryMemory{}
	basicUserID, err := basicUserRepository.AddBasicUser(ctx, channelID, basicUser)
	require.NoError(t, err)

	err = basicUser.SetUUID(basicUserID)
	require.NoError(t, err)

	actorUser := actor.Actor{BasicUser: basicUser}

	clock := mocks.FixedClock{}
	repo := memory.NewIncidentRepositoryMemory(clock, basicUserRepository)

	svc := NewIncidentService(repo)

	// CreateIncident
	params1 := api.CreateIncidentParams{
		Number:           "ABC123",
		ShortDescription: "Some incident 1",
		Description:      "Nice one",
	}
	inc1ID, err := svc.CreateIncident(ctx, channelID, actorUser, params1)
	require.NoError(t, err)

	params2 := api.CreateIncidentParams{
		Number:           "DEF456",
		ShortDescription: "Some incident 2",
	}
	inc2ID, err := svc.CreateIncident(ctx, channelID, actorUser, params2)
	require.NoError(t, err)

	// ListIncidents
	paginationParams := new(mocks.PaginationParamsMock)
	paginationParams.On("Page").Return(uint(1))
	paginationParams.On("ItemsPerPage").Return(uint(10))

	list, err := svc.ListIncidents(ctx, channelID, actorUser, paginationParams)
	require.NoError(t, err)

	incidents := list.Result
	assert.Len(t, incidents, 2)

	assert.Equal(t, inc1ID, incidents[0].UUID())
	assert.Equal(t, inc2ID, incidents[1].UUID())

	// GetIncident
	retInc1, err := svc.GetIncident(ctx, channelID, actorUser, inc1ID)
	require.NoError(t, err)
	assert.Equal(t, retInc1.Number, params1.Number)
	assert.Equal(t, retInc1.ShortDescription, params1.ShortDescription)
	assert.Equal(t, retInc1.UUID(), inc1ID)

	retInc2, err := svc.GetIncident(ctx, channelID, actorUser, inc2ID)
	require.NoError(t, err)
	assert.Equal(t, retInc2.Number, params2.Number)
	assert.Equal(t, retInc2.ShortDescription, params2.ShortDescription)
	assert.Equal(t, retInc2.UUID(), inc2ID)
}

func Test_incidentService_UpdateIncident(t *testing.T) {
	channelID := ref.ChannelID("e27ddcd0-0e1f-4bc5-93df-f6f04155beec")
	ctx := context.Background()

	basicUser := user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}

	basicUserRepository := &memory.BasicUserRepositoryMemory{}
	basicUserID, err := basicUserRepository.AddBasicUser(ctx, channelID, basicUser)
	require.NoError(t, err)

	err = basicUser.SetUUID(basicUserID)
	require.NoError(t, err)

	actorUser := actor.Actor{BasicUser: basicUser}

	clock := mocks.FixedClock{}
	repo := memory.NewIncidentRepositoryMemory(clock, basicUserRepository)

	svc := NewIncidentService(repo)

	// CreateIncident
	params := api.CreateIncidentParams{
		Number:           "ABC123",
		ExternalID:       "111222333",
		ShortDescription: "Some incident 1",
		Description:      "Nice one",
	}
	incID, err := svc.CreateIncident(ctx, channelID, actorUser, params)
	require.NoError(t, err)

	// GetIncident
	origInc, err := svc.GetIncident(ctx, channelID, actorUser, incID)
	require.NoError(t, err)

	updateParams := api.UpdateIncidentParams{
		ShortDescription: "Some updated short description",
	}
	updatedIncID, err := svc.UpdateIncident(ctx, channelID, actorUser, incID, updateParams)
	require.NoError(t, err)

	assert.Equal(t, updatedIncID, incID)

	updatedInc, err := svc.GetIncident(ctx, channelID, actorUser, incID)
	require.NoError(t, err)
	assert.Equal(t, origInc.Number, updatedInc.Number)
	assert.Equal(t, origInc.ExternalID, updatedInc.ExternalID)
	assert.Equal(t, updateParams.ShortDescription, updatedInc.ShortDescription)
	assert.Equal(t, "", updatedInc.Description)
	assert.Equal(t, updatedInc.UUID(), incID)
}
