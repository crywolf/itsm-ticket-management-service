package incidentsvc

import (
	"context"
	"testing"
	"time"

	fieldengineer "github.com/crywolf/itsm-ticket-management-service/internal/domain/field_engineer"
	fieldengineersvc "github.com/crywolf/itsm-ticket-management-service/internal/domain/field_engineer/service"
	tsession "github.com/crywolf/itsm-ticket-management-service/internal/domain/field_engineer/time_session"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/incident"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/crywolf/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/crywolf/itsm-ticket-management-service/internal/mocks"
	"github.com/crywolf/itsm-ticket-management-service/internal/repository/memory"
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

	clock := mocks.NewFixedClock()
	fieldEngineerRepository := memory.NewFieldEngineerRepositoryMemory(clock, basicUserRepository)
	incidentRepository := memory.NewIncidentRepositoryMemory(clock, basicUserRepository, fieldEngineerRepository)

	svc := NewIncidentService(incidentRepository, fieldEngineerRepository)

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

	fieldEngineer := fieldengineer.FieldEngineer{
		BasicUser: basicUser,
	}
	err = fieldEngineer.CreatedUpdated.SetCreatedBy(basicUser)
	require.NoError(t, err)
	err = fieldEngineer.CreatedUpdated.SetUpdatedBy(basicUser)
	require.NoError(t, err)

	clock := mocks.NewFixedClock()
	fieldEngineerRepository := memory.NewFieldEngineerRepositoryMemory(clock, basicUserRepository)
	fieldEngineerID, err := fieldEngineerRepository.AddFieldEngineer(ctx, channelID, fieldEngineer)
	require.NoError(t, err)
	err = fieldEngineer.SetUUID(fieldEngineerID)
	require.NoError(t, err)

	incidentRepository := memory.NewIncidentRepositoryMemory(clock, basicUserRepository, fieldEngineerRepository)

	svc := NewIncidentService(incidentRepository, fieldEngineerRepository)

	feUUID := api.UUID(fieldEngineer.UUID().String())
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

	// add some time
	clock.AddTime(100 * time.Second)

	// trying update with non-existing field engineer
	wrongFeUUID := api.UUID("3d334abe-f289-42a5-9742-72c3133768c2")
	updateParamsNonExistingFE := api.UpdateIncidentParams{
		ShortDescription: "Some updated short description",
		FieldEngineerID:  &wrongFeUUID,
	}
	_, err = svc.UpdateIncident(ctx, channelID, actorUser, incID, updateParamsNonExistingFE)
	// it should return error
	require.Error(t, err)
	assert.EqualError(t, err, "cannot assign field engineer: error loading field engineer from repository: record was not found")

	// update with existing FE
	updateParams := api.UpdateIncidentParams{
		ShortDescription: "Some updated short description",
		FieldEngineerID:  &feUUID,
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
	assert.NotNil(t, updatedInc.FieldEngineerID)
	assert.Equal(t, fieldEngineer.UUID(), *updatedInc.FieldEngineerID)

	// timestamps
	assert.Equal(t, origInc.CreatedUpdated.CreatedBy(), updatedInc.CreatedUpdated.CreatedBy())
	assert.Equal(t, origInc.CreatedUpdated.CreatedAt(), updatedInc.CreatedUpdated.CreatedAt())
	assert.Equal(t, origInc.CreatedUpdated.UpdatedBy(), updatedInc.CreatedUpdated.UpdatedBy())
	// timestamp updatedAt should change
	assert.NotEqual(t, origInc.CreatedUpdated.UpdatedAt(), updatedInc.CreatedUpdated.UpdatedAt())
	assert.Equal(t, updatedInc.UUID(), incID)
}

func Test_incidentService_StartWorking_and_StopWorking(t *testing.T) {
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

	clock := mocks.NewFixedClock()
	fieldEngineerRepository := memory.NewFieldEngineerRepositoryMemory(clock, basicUserRepository)
	feSvc := fieldengineersvc.NewFieldEngineerService(fieldEngineerRepository)

	incidentRepository := memory.NewIncidentRepositoryMemory(clock, basicUserRepository, fieldEngineerRepository)
	svc := NewIncidentService(incidentRepository, fieldEngineerRepository)

	// create field engineer
	fe := fieldengineer.FieldEngineer{BasicUser: basicUser}
	err = fe.CreatedUpdated.SetCreatedBy(basicUser)
	require.NoError(t, err)
	err = fe.CreatedUpdated.SetUpdatedBy(basicUser)
	require.NoError(t, err)
	feID, err := fieldEngineerRepository.AddFieldEngineer(ctx, channelID, fe)
	require.NoError(t, err)

	// set actor as field engineer
	actorUser.SetFieldEngineerID(&feID)

	// CreateIncident
	feUUID := api.UUID(feID)
	incParams := api.CreateIncidentParams{
		Number:           "ABC123",
		ShortDescription: "Some incident 1",
		Description:      "Nice one",
		FieldEngineerID:  &feUUID,
	}
	incID, err := svc.CreateIncident(ctx, channelID, actorUser, incParams)
	require.NoError(t, err)

	// StartWorking
	remote := true
	err = svc.StartWorking(ctx, channelID, actorUser, incID, api.IncidentStartWorkingParams{Remote: remote}, clock)
	require.NoError(t, err)

	// GetIncident
	updatedInc, err := svc.GetIncident(ctx, channelID, actorUser, incID)
	require.NoError(t, err)

	assert.Equal(t, incParams.Number, updatedInc.Number)
	assert.Equal(t, incParams.ExternalID, updatedInc.ExternalID)
	assert.Equal(t, incident.StateInProgress, updatedInc.State())
	assert.Len(t, updatedInc.Timelogs, 1)
	assert.True(t, updatedInc.HasOpenTimelog())
	assert.NotEmpty(t, updatedInc.OpenTimelog().Start)
	assert.Empty(t, updatedInc.OpenTimelog().End)
	assert.Empty(t, updatedInc.OpenTimelog().Work)
	assert.Equal(t, remote, updatedInc.OpenTimelog().Remote)

	// GetFieldEngineer
	updatedFe, err := feSvc.GetFieldEngineer(ctx, channelID, actorUser, feID)
	require.NoError(t, err)

	assert.Len(t, updatedFe.TimeSessions, 1)
	assert.Equal(t, true, updatedFe.HasOpenTimeSession())
	assert.NotNil(t, updatedFe.OpenTimeSession())
	openTS := updatedFe.OpenTimeSession()
	assert.Equal(t, tsession.StateWork, openTS.State())
	assert.Len(t, openTS.Incidents, 1)
	assert.Equal(t, incID, openTS.Incidents[0].IncidentID)

	clock.AddTime(2 * time.Hour)
	err = svc.StopWorking(ctx, channelID, actorUser, incID, api.IncidentStopWorkingParams{VisitSummary: "some message"}, clock)
	require.NoError(t, err)

	// GetIncident
	updatedInc, err = svc.GetIncident(ctx, channelID, actorUser, incID)
	require.NoError(t, err)

	assert.Equal(t, incParams.Number, updatedInc.Number)
	assert.Equal(t, incParams.ExternalID, updatedInc.ExternalID)
	assert.Equal(t, incident.StateInProgress, updatedInc.State())
	assert.Len(t, updatedInc.Timelogs, 1)
	assert.False(t, updatedInc.HasOpenTimelog())
	assert.Nil(t, updatedInc.OpenTimelog())

	// GetIncidentTimelog
	timelogID := updatedInc.Timelogs[0]
	timelog, err := svc.GetIncidentTimelog(ctx, channelID, actorUser, incID, timelogID)
	require.NoError(t, err)

	assert.NotEmpty(t, timelog.End)
	assert.Equal(t, clock.NowFormatted(), timelog.End)
	assert.NotEmpty(t, timelog.Work)
}
