package rest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/mocks"
	"github.com/KompiTech/itsm-ticket-management-service/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateIncidentHandler(t *testing.T) {
}

func TestGetIncidentHandler(t *testing.T) {
	logger, _ := testutils.NewTestLogger()
	defer func() { _ = logger.Sync() }()

	channelID := "e27ddcd0-0e1f-4bc5-93df-f6f04155beec"
	bearerToken := "some valid Bearer token"

	t.Run("when incident exists", func(t *testing.T) {
		uuid := "cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0"
		retInc := incident.Incident{
			ShortDescription: "Test incident 1",
		}
		err := retInc.SetUUID(ref.UUID(uuid))
		require.NoError(t, err)
		state, err := incident.NewStateFromString("new")
		require.NoError(t, err)
		err = retInc.SetState(state)
		require.NoError(t, err)
		err = retInc.CreatedUpdated.SetCreated("8540d943-8ccd-4ff1-8a08-0c3aa338c58e", "2021-04-01T12:34:56+02:00")
		require.NoError(t, err)
		err = retInc.CreatedUpdated.SetUpdated("8540d943-8ccd-4ff1-8a08-0c3aa338c58e", "2021-04-01T12:34:56+02:00")
		require.NoError(t, err)

		incidentSvc := new(mocks.IncidentServiceMock)
		incidentSvc.On("GetIncident", ref.UUID(uuid), ref.ChannelID(channelID)).
			Return(retInc, nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
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

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		expectedJSON := `{
			"uuid":"cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0",
			"short_description":"Test incident 1",
			"state":"new",
			"created_by":"8540d943-8ccd-4ff1-8a08-0c3aa338c58e",
			"created_at":"2021-04-01T12:34:56+02:00",
			"updated_by":"8540d943-8ccd-4ff1-8a08-0c3aa338c58e",
			"updated_at":"2021-04-01T12:34:56+02:00",
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

	t.Run("when some incidents were found", func(t *testing.T) {
		expectedJSON := `{
			"result":
				[{
					"uuid":"cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0",
					"short_description":"Test incident 1",
					"state":"new",
					"created_by":"8540d943-8ccd-4ff1-8a08-0c3aa338c58e",
					"created_at":"2021-04-01T12:34:56+02:00",
					"updated_by":"8540d943-8ccd-4ff1-8a08-0c3aa338c58e",
					"updated_at":"2021-04-01T12:34:56+02:00",
					"_links":{
						"self":{"href":"http://service.url/incidents/cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0"},
						"CancelIncident":{"href":"http://service.url/incidents/cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0/cancel"},
						"IncidentStartWorking":{"href":"http://service.url/incidents/cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0/start_working"}
					}
				},
				{
					"uuid": "0ac5ebce-17e7-4edc-9552-fefe16e127fb",
					"short_description":"Test incident 2",
					"state":"resolved",
					"created_by":"8540d943-8ccd-4ff1-8a08-0c3aa338c58e",
					"created_at": "2021-04-11T00:45:42+02:00",
					"updated_by":"cd00bc0a-cc45-498c-9d2c-4d7e52efcd30",
					"updated_at":"2021-04-02T09:10:32+02:00",
					"_links":{
						"self":{"href":"http://service.url/incidents/0ac5ebce-17e7-4edc-9552-fefe16e127fb"},
						"CancelIncident":{"href":"http://service.url/incidents/0ac5ebce-17e7-4edc-9552-fefe16e127fb/cancel"},
						"IncidentStartWorking":{"href":"http://service.url/incidents/0ac5ebce-17e7-4edc-9552-fefe16e127fb/start_working"}
					}
				}],
			"_links":{
				"self":{"href":"http://service.url/incidents"}
			}
		}`

		var list []incident.Incident

		fInc1 := incident.Incident{
			ShortDescription: "Test incident 1",
		}
		err := fInc1.SetUUID("cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0")
		require.NoError(t, err)
		state, err := incident.NewStateFromString("new")
		require.NoError(t, err)
		err = fInc1.SetState(state)
		require.NoError(t, err)
		err = fInc1.CreatedUpdated.SetCreated("8540d943-8ccd-4ff1-8a08-0c3aa338c58e", "2021-04-01T12:34:56+02:00")
		require.NoError(t, err)
		err = fInc1.CreatedUpdated.SetUpdated("8540d943-8ccd-4ff1-8a08-0c3aa338c58e", "2021-04-01T12:34:56+02:00")
		require.NoError(t, err)
		list = append(list, fInc1)

		fInc2 := incident.Incident{
			ShortDescription: "Test incident 2",
		}
		err = fInc2.SetUUID("0ac5ebce-17e7-4edc-9552-fefe16e127fb")
		require.NoError(t, err)
		state, err = incident.NewStateFromString("resolved")
		require.NoError(t, err)
		err = fInc2.SetState(state)
		require.NoError(t, err)
		err = fInc2.CreatedUpdated.SetCreated("8540d943-8ccd-4ff1-8a08-0c3aa338c58e", "2021-04-11T00:45:42+02:00")
		require.NoError(t, err)
		err = fInc2.CreatedUpdated.SetUpdated("cd00bc0a-cc45-498c-9d2c-4d7e52efcd30", "2021-04-02T09:10:32+02:00")
		require.NoError(t, err)
		list = append(list, fInc2)

		incidentSvc := new(mocks.IncidentServiceMock)
		incidentSvc.On("ListIncidents", ref.ChannelID(channelID)).Return(list, nil)

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
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

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})
}
