package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crywolf/itsm-ticket-management-service/internal/domain"
	fieldengineer "github.com/crywolf/itsm-ticket-management-service/internal/domain/field_engineer"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/incident"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/crywolf/itsm-ticket-management-service/internal/mocks"
	"github.com/crywolf/itsm-ticket-management-service/internal/repository"
	"github.com/crywolf/itsm-ticket-management-service/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateIncidentHandler(t *testing.T) {
	logger, _ := testutils.NewTestLogger()
	defer func() { _ = logger.Sync() }()

	channelID := "e27ddcd0-0e1f-4bc5-93df-f6f04155beec"
	bearerToken := "some valid Bearer token"

	actorUser := actor.Actor{
		BasicUser: user.BasicUser{
			ExternalUserUUID: "cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0",
		},
	}
	err := actorUser.BasicUser.SetUUID("8183eaca-56c0-41d9-9291-1d295dd53763")
	require.NoError(t, err)

	t.Parallel()

	t.Run("when body payload is not valid JSON", func(t *testing.T) {
		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalLocationAddress: "http://service.url",
			ExternalUserService:     us,
		})

		payload := []byte(`{"invalid json request"}`)

		body := bytes.NewReader(payload)
		req := httptest.NewRequest("POST", "/incidents", body)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		defer func() { _ = resp.Body.Close() }()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		us.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		expectedJSON := `{"error":"Request body contains badly-formed JSON (at position 24): invalid character '}' after object key"}`
		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})

	t.Run("when body payload is not valid (ie. validation fails)", func(t *testing.T) {
		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalLocationAddress: "http://service.url",
			ExternalUserService:     us,
		})

		payload := []byte(`{
			"field_engineer": null,
			"description": "incident with required fields missing"
		}`)

		body := bytes.NewReader(payload)
		req := httptest.NewRequest("POST", "/incidents", body)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		defer func() { _ = resp.Body.Close() }()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		us.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		expectedJSON := `{"error":"'number' is a required field, 'short_description' is a required field"}`
		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})

	t.Run("when body payload is valid", func(t *testing.T) {
		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		incidentSvc := new(mocks.IncidentServiceMock)
		incidentSvc.On("CreateIncident", ref.ChannelID(channelID), actorUser, mock.AnythingOfType("api.CreateIncidentParams")).
			Return(ref.UUID("38316161-3035-4864-ad30-6231392d3433"), nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			IncidentService:         incidentSvc,
			ExternalLocationAddress: "http://service.url",
			ExternalUserService:     us,
		})

		payload := []byte(`{
			"number": "INC123",
			"short_description": "some test incident 1"
		}`)

		body := bytes.NewReader(payload)
		req := httptest.NewRequest("POST", "/incidents", body)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		us.AssertExpectations(t)
		incidentSvc.AssertExpectations(t)

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Status code")
		expectedLocation := "http://service.url/incidents/38316161-3035-4864-ad30-6231392d3433"
		assert.Equal(t, expectedLocation, resp.Header.Get("Location"), "Location header")
	})
}

func TestUpdateIncidentHandler(t *testing.T) {
	logger, _ := testutils.NewTestLogger()
	defer func() { _ = logger.Sync() }()

	channelID := "e27ddcd0-0e1f-4bc5-93df-f6f04155beec"
	bearerToken := "some valid Bearer token"

	actorUser := actor.Actor{
		BasicUser: user.BasicUser{
			ExternalUserUUID: "cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0",
		},
	}
	err := actorUser.BasicUser.SetUUID("8183eaca-56c0-41d9-9291-1d295dd53763")
	require.NoError(t, err)

	t.Parallel()

	t.Run("when body payload is not valid JSON", func(t *testing.T) {
		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalLocationAddress: "http://service.url",
			ExternalUserService:     us,
		})

		payload := []byte(`{"invalid json request"}`)

		body := bytes.NewReader(payload)
		req := httptest.NewRequest("PATCH", "/incidents/7e0d38d1-e5f5-4211-b2aa-3b142e4da80e", body)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		defer func() { _ = resp.Body.Close() }()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		us.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		expectedJSON := `{"error":"Request body contains badly-formed JSON (at position 24): invalid character '}' after object key"}`
		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})

	t.Run("when body payload is not valid (ie. validation fails)", func(t *testing.T) {
		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalLocationAddress: "http://service.url",
			ExternalUserService:     us,
		})

		payload := []byte(`{
			"description": "changing the description",
			"field_engineer": ""
		}`)

		body := bytes.NewReader(payload)
		req := httptest.NewRequest("PATCH", "/incidents/7e0d38d1-e5f5-4211-b2aa-3b142e4da80e", body)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		defer func() { _ = resp.Body.Close() }()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		us.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		expectedJSON := `{"error":"'short_description' is a required field, 'field_engineer' must be a valid version 4 UUID"}`
		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})

	t.Run("when body payload is valid", func(t *testing.T) {
		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		incidentSvc := new(mocks.IncidentServiceMock)
		incidentSvc.On("UpdateIncident", ref.ChannelID(channelID), actorUser, ref.UUID("7e0d38d1-e5f5-4211-b2aa-3b142e4da80e"),
			mock.AnythingOfType("api.UpdateIncidentParams")).Return(ref.UUID("7e0d38d1-e5f5-4211-b2aa-3b142e4da80e"), nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			IncidentService:         incidentSvc,
			ExternalLocationAddress: "http://service.url",
			ExternalUserService:     us,
		})

		payload := []byte(`{
			"short_description": "changed description"
		}`)

		body := bytes.NewReader(payload)
		req := httptest.NewRequest("PATCH", "/incidents/7e0d38d1-e5f5-4211-b2aa-3b142e4da80e", body)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		us.AssertExpectations(t)
		incidentSvc.AssertExpectations(t)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode, "Status code")
		expectedLocation := "http://service.url/incidents/7e0d38d1-e5f5-4211-b2aa-3b142e4da80e"
		assert.Equal(t, expectedLocation, resp.Header.Get("Location"), "Location header")
	})
}

func TestGetIncidentHandler(t *testing.T) {
	logger, _ := testutils.NewTestLogger()
	defer func() { _ = logger.Sync() }()

	channelID := "e27ddcd0-0e1f-4bc5-93df-f6f04155beec"
	bearerToken := "some valid Bearer token"

	createdByUser := user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	err := createdByUser.SetUUID("cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0")
	require.NoError(t, err)

	fieldEngineer := &fieldengineer.FieldEngineer{
		BasicUser: user.BasicUser{
			ExternalUserUUID: "5d5ef779-17cb-413a-aa4b-7bc0a80bf230",
			Name:             "Alois",
			Surname:          "Vomacka",
			OrgDisplayName:   "CGI",
			OrgName:          "1233ae78-cb08-4fd3-9d59-b3b8b07e08fc.kompitech.com",
		},
	}
	err = fieldEngineer.SetUUID("1adb8393-cff0-489c-a82f-3fe5d15708d4")
	require.NoError(t, err)
	fieldEngineerUUID := fieldEngineer.UUID()

	actorUser := actor.Actor{
		BasicUser: fieldEngineer.BasicUser,
	}
	actorUser.SetFieldEngineerID(&fieldEngineerUUID)

	t.Parallel()

	t.Run("when 'channel-id' header is missing", func(t *testing.T) {
		uuid := "cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0"

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalLocationAddress: "http://service.url",
		})

		req := httptest.NewRequest("GET", "/incidents/"+uuid, nil)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		defer func() { _ = resp.Body.Close() }()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		expectedJSON := `{"error":"'channel-id' header missing or invalid"}`
		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})

	t.Run("when incident does not exist", func(t *testing.T) {
		uuid := "cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0"

		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		incidentSvc := new(mocks.IncidentServiceMock)
		incidentSvc.On("GetIncident", ref.ChannelID(channelID), actorUser, ref.UUID(uuid)).
			Return(incident.Incident{}, domain.NewErrorf(domain.ErrorCodeNotFound, "error from repository"))

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalUserService:     us,
			IncidentService:         incidentSvc,
			ExternalLocationAddress: "http://service.url",
		})

		req := httptest.NewRequest("GET", "/incidents/"+uuid, nil)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		defer func() { _ = resp.Body.Close() }()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		us.AssertExpectations(t)
		incidentSvc.AssertExpectations(t)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		expectedJSON := `{"error":"incident not found"}`
		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})

	t.Run("when incident exists", func(t *testing.T) {
		uuid := "cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0"
		retInc := incident.Incident{
			Number:           "A123456",
			ShortDescription: "Test incident 1",
			FieldEngineerID:  &fieldEngineerUUID,
		}
		err := retInc.SetUUID(ref.UUID(uuid))
		require.NoError(t, err)
		state, err := incident.NewStateFromString("new")
		require.NoError(t, err)
		err = retInc.SetState(state)
		require.NoError(t, err)
		err = retInc.CreatedUpdated.SetCreated(createdByUser, "2021-04-01T12:34:56+02:00")
		require.NoError(t, err)
		err = retInc.CreatedUpdated.SetUpdated(createdByUser, "2021-04-01T12:34:56+02:00")
		require.NoError(t, err)

		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		incidentSvc := new(mocks.IncidentServiceMock)
		incidentSvc.On("GetIncident", ref.ChannelID(channelID), actorUser, ref.UUID(uuid)).
			Return(retInc, nil)

		feSvc := new(mocks.FieldEngineerServiceMock)
		feSvc.On("GetFieldEngineer", ref.ChannelID(channelID), actorUser, fieldEngineerUUID).
			Return(*fieldEngineer, nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalUserService:     us,
			IncidentService:         incidentSvc,
			FieldEngineerService:    feSvc,
			ExternalLocationAddress: "http://service.url",
		})

		req := httptest.NewRequest("GET", "/incidents/"+uuid, nil)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		defer func() { _ = resp.Body.Close() }()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		us.AssertExpectations(t)
		incidentSvc.AssertExpectations(t)
		feSvc.AssertExpectations(t)

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		expectedJSON := `{
			"uuid":"cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0",
			"number": "A123456",
			"short_description":"Test incident 1",
			"field_engineer":"1adb8393-cff0-489c-a82f-3fe5d15708d4",
			"state":"new",
			"created_by":"cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0",
			"created_at":"2021-04-01T12:34:56+02:00",
			"updated_by":"cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0",
			"updated_at":"2021-04-01T12:34:56+02:00",
			"_embedded":{
				"field_engineer":{
					"_links": {
						"self": {"href": "http://service.url/field_engineers/1adb8393-cff0-489c-a82f-3fe5d15708d4"}
					},
				    "external_user_uuid": "5d5ef779-17cb-413a-aa4b-7bc0a80bf230",
					"name":"Alois",
					"surname":"Vomacka",
					"org_name":"1233ae78-cb08-4fd3-9d59-b3b8b07e08fc.kompitech.com",
					"org_display_name":"CGI",
					"uuid":"1adb8393-cff0-489c-a82f-3fe5d15708d4"
				},
				"created_by":{
					"_links": {
						"self": {"href": "http://service.url/basic_users/cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0"}
					},
				    "external_user_uuid": "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
					"name":"Alfred",
					"surname":"Koletschko",
					"org_name":"a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
					"org_display_name":"KompiTech",
				    "uuid": "cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0"
				}
			},
			"_links":{
				"self":{"href":"http://service.url/incidents/cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0"},
				"CancelIncident":{"href":"http://service.url/incidents/cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0/cancel"},
				"IncidentStartWorking":{"href":"http://service.url/incidents/cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0/start_working"}
			}
		}`
		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})
}

func TestListIncidentsHandler(t *testing.T) {
	logger, _ := testutils.NewTestLogger()
	defer func() { _ = logger.Sync() }()

	channelID := "e27ddcd0-0e1f-4bc5-93df-f6f04155beec"
	bearerToken := "some valid Bearer token"

	createdByUser := user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	err := createdByUser.SetUUID("8183eaca-56c0-41d9-9291-1d295dd53763")
	require.NoError(t, err)

	fieldEngineer := &fieldengineer.FieldEngineer{
		BasicUser: user.BasicUser{
			ExternalUserUUID: "5d5ef779-17cb-413a-aa4b-7bc0a80bf230",
			Name:             "Alois",
			Surname:          "Vomacka",
			OrgDisplayName:   "CGI",
			OrgName:          "1233ae78-cb08-4fd3-9d59-b3b8b07e08fc.kompitech.com",
		},
	}
	err = fieldEngineer.SetUUID("1adb8393-cff0-489c-a82f-3fe5d15708d4")
	require.NoError(t, err)
	fieldEngineerUUID := fieldEngineer.UUID()

	actorUser := actor.Actor{
		BasicUser: fieldEngineer.BasicUser,
	}
	actorUser.SetFieldEngineerID(&fieldEngineerUUID)

	t.Parallel()

	t.Run("when no incidents were found", func(t *testing.T) {
		expectedJSON := `{
			"total":0,
			"size":0,
			"page":1,
			"_links":{
				"self":{"href":"http://service.url/incidents"},
				"first":{"href":"http://service.url/incidents"},
				"last":{"href":"http://service.url/incidents"}
			}
		}`

		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		incidentSvc := new(mocks.IncidentServiceMock)
		var emptyList []incident.Incident
		result := repository.IncidentList{
			Result: emptyList,
			Pagination: &repository.Pagination{
				Total: 0,
				Size:  0,
				Page:  1,
				First: 1,
				Last:  1,
			},
		}
		incidentSvc.On("ListIncidents", ref.ChannelID(channelID), actorUser, mock.AnythingOfType("*converters.paginationParams")).Return(result, nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalUserService:     us,
			IncidentService:         incidentSvc,
			ExternalLocationAddress: "http://service.url",
		})

		req := httptest.NewRequest("GET", "/incidents", nil)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		defer func() { _ = resp.Body.Close() }()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		us.AssertExpectations(t)
		incidentSvc.AssertExpectations(t)

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})

	t.Run("when some incidents were found", func(t *testing.T) {
		expectedJSON := `{
			"total":2,
			"size":2,
			"page":1,
			"_embedded":
				[{
					"uuid":"cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0",
					"number": "Accc265871",
					"short_description":"Test incident 1",
					"field_engineer":"1adb8393-cff0-489c-a82f-3fe5d15708d4",
					"state":"new",
					"created_by":"8183eaca-56c0-41d9-9291-1d295dd53763",
					"created_at":"2021-04-01T12:34:56+02:00",
					"updated_by":"8183eaca-56c0-41d9-9291-1d295dd53763",
					"updated_at":"2021-04-01T12:34:56+02:00",
					"_embedded":{
						"field_engineer":{
							"_links": {
								"self": {"href": "http://service.url/field_engineers/1adb8393-cff0-489c-a82f-3fe5d15708d4"}
							},
							"external_user_uuid": "5d5ef779-17cb-413a-aa4b-7bc0a80bf230",
							"name":"Alois",
							"surname":"Vomacka",
							"org_name":"1233ae78-cb08-4fd3-9d59-b3b8b07e08fc.kompitech.com",
							"org_display_name":"CGI",
							"uuid":"1adb8393-cff0-489c-a82f-3fe5d15708d4"
						}
					},
					"_links":{
						"self":{"href":"http://service.url/incidents/cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0"},
						"CancelIncident":{"href":"http://service.url/incidents/cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0/cancel"},
						"IncidentStartWorking":{"href":"http://service.url/incidents/cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0/start_working"}
					}
				},
				{
					"uuid": "0ac5ebce-17e7-4edc-9552-fefe16e127fb",
					"number": "555555",
					"short_description":"Test incident 2 - with field engineer assigned",
					"field_engineer":null,
					"state":"resolved",
					"created_by":"8183eaca-56c0-41d9-9291-1d295dd53763",
					"created_at": "2021-04-11T00:45:42+02:00",
					"updated_by":"8183eaca-56c0-41d9-9291-1d295dd53763",
					"updated_at":"2021-04-02T09:10:32+02:00",
					"_links":{
						"self":{"href":"http://service.url/incidents/0ac5ebce-17e7-4edc-9552-fefe16e127fb"}
					}
				}],
			"_links":{
				"self":{"href":"http://service.url/incidents"},
				"first":{"href":"http://service.url/incidents"},
				"last":{"href":"http://service.url/incidents"}
			}
		}`

		var list []incident.Incident

		basicUser2 := user.BasicUser{
			ExternalUserUUID: "cd00bc0a-cc45-498c-9d2c-4d7e52efcd30",
		}
		_ = basicUser2.SetUUID("8183eaca-56c0-41d9-9291-1d295dd53763")

		fInc1 := incident.Incident{
			Number:           "Accc265871",
			ShortDescription: "Test incident 1",
			FieldEngineerID:  &fieldEngineerUUID,
		}
		err := fInc1.SetUUID("cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0")
		require.NoError(t, err)
		state, err := incident.NewStateFromString("new")
		require.NoError(t, err)
		err = fInc1.SetState(state)
		require.NoError(t, err)
		err = fInc1.CreatedUpdated.SetCreated(createdByUser, "2021-04-01T12:34:56+02:00")
		require.NoError(t, err)
		err = fInc1.CreatedUpdated.SetUpdated(createdByUser, "2021-04-01T12:34:56+02:00")
		require.NoError(t, err)
		list = append(list, fInc1)

		fInc2 := incident.Incident{
			Number:           "555555",
			ShortDescription: "Test incident 2 - with field engineer assigned",
		}
		err = fInc2.SetUUID("0ac5ebce-17e7-4edc-9552-fefe16e127fb")
		require.NoError(t, err)
		state, err = incident.NewStateFromString("resolved")
		require.NoError(t, err)
		err = fInc2.SetState(state)
		require.NoError(t, err)
		err = fInc2.CreatedUpdated.SetCreated(createdByUser, "2021-04-11T00:45:42+02:00")
		require.NoError(t, err)
		err = fInc2.CreatedUpdated.SetUpdated(basicUser2, "2021-04-02T09:10:32+02:00")
		require.NoError(t, err)
		list = append(list, fInc2)

		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		incidentSvc := new(mocks.IncidentServiceMock)
		result := repository.IncidentList{
			Result: list,
			Pagination: &repository.Pagination{
				Total: 2,
				Size:  2,
				Page:  1,
				First: 1,
				Last:  1,
			},
		}
		incidentSvc.On("ListIncidents", ref.ChannelID(channelID), actorUser, mock.AnythingOfType("*converters.paginationParams")).Return(result, nil)

		feSvc := new(mocks.FieldEngineerServiceMock)
		feSvc.On("GetFieldEngineer", ref.ChannelID(channelID), actorUser, fieldEngineerUUID).
			Return(*fieldEngineer, nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalUserService:     us,
			IncidentService:         incidentSvc,
			FieldEngineerService:    feSvc,
			ExternalLocationAddress: "http://service.url",
		})

		req := httptest.NewRequest("GET", "/incidents", nil)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		defer func() { _ = resp.Body.Close() }()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		us.AssertExpectations(t)
		incidentSvc.AssertExpectations(t)
		feSvc.AssertExpectations(t)

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})

	t.Run("when 'page' parameter in HTTP query is incorrect", func(t *testing.T) {
		expectedJSON := `{"error":"incorrect 'page' parameter: '0'"}`

		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalUserService:     us,
			ExternalLocationAddress: "http://service.url",
		})

		req := httptest.NewRequest("GET", "/incidents?page=0", nil)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		defer func() { _ = resp.Body.Close() }()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		us.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})
}

func TestIncidentStartWorkingHandler(t *testing.T) {
	logger, _ := testutils.NewTestLogger()
	defer func() { _ = logger.Sync() }()

	channelID := "e27ddcd0-0e1f-4bc5-93df-f6f04155beec"
	bearerToken := "some valid Bearer token"

	fieldEngineer := &fieldengineer.FieldEngineer{
		BasicUser: user.BasicUser{
			ExternalUserUUID: "5d5ef779-17cb-413a-aa4b-7bc0a80bf230",
			Name:             "Alois",
			Surname:          "Vomacka",
			OrgDisplayName:   "CGI",
			OrgName:          "1233ae78-cb08-4fd3-9d59-b3b8b07e08fc.kompitech.com",
		},
	}
	err := fieldEngineer.SetUUID("1adb8393-cff0-489c-a82f-3fe5d15708d4")
	require.NoError(t, err)
	//	fieldEngineerUUID := fieldEngineer.UUID()

	actorUser := actor.Actor{
		BasicUser: fieldEngineer.BasicUser,
	}
	//	actorUser.SetFieldEngineerID(&fieldEngineerUUID)

	t.Parallel()

	t.Run("everything is ok", func(t *testing.T) {
		uuid := "cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0"
		retInc := incident.Incident{
			Number:           "A123456",
			ShortDescription: "Test incident 1",
			//			FieldEngineerID:  &fieldEngineerUUID,
		}
		err := retInc.SetUUID(ref.UUID(uuid))
		require.NoError(t, err)

		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actorUser, nil)

		incidentSvc := new(mocks.IncidentServiceMock)
		incidentSvc.On("StartWorking", ref.ChannelID(channelID), actorUser, ref.UUID(uuid), mock.AnythingOfType("api.IncidentStartWorkingParams")).
			Return(nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalUserService:     us,
			IncidentService:         incidentSvc,
			ExternalLocationAddress: "http://service.url",
		})

		payload := []byte(`{}`)

		body := bytes.NewReader(payload)
		req := httptest.NewRequest("POST", "/incidents/"+uuid+"/start_working", body)
		req.Header.Set("channel-id", channelID)
		req.Header.Set("authorization", bearerToken)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		us.AssertExpectations(t)
		incidentSvc.AssertExpectations(t)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode, "Status code")
		expectedLocation := "http://service.url/incidents/cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0"
		assert.Equal(t, expectedLocation, resp.Header.Get("Location"), "Location header")
	})
}
