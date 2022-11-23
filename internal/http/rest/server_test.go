package rest

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crywolf/itsm-ticket-management-service/internal/domain"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/crywolf/itsm-ticket-management-service/internal/mocks"
	"github.com/crywolf/itsm-ticket-management-service/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	logger, _ := testutils.NewTestLogger()
	defer func() { _ = logger.Sync() }()

	channelID := "e27ddcd0-0e1f-4bc5-93df-f6f04155beec"
	bearerToken := "some valid Bearer token"

	t.Parallel()

	t.Run("when 'authorization' header with Bearer token is missing", func(t *testing.T) {
		server := NewServer(Config{
			Addr:   "service.url",
			Logger: logger,
		})

		req := httptest.NewRequest("GET", "/someEndpoint", nil)

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

		expectedJSON := `{"error":"'authorization' header missing or invalid"}`
		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})

	t.Run("when 'channel-id' header is missing", func(t *testing.T) {
		server := NewServer(Config{
			Addr:   "service.url",
			Logger: logger,
		})

		req := httptest.NewRequest("GET", "/someEndpoint", nil)
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

	t.Run("when user service failed to retrieve Actor and put it in the request context (some GRPC client error)", func(t *testing.T) {
		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actor.Actor{}, domain.WrapErrorf(errors.New("some user service GRPC error"), domain.ErrorCodeUnknown, "authorization failed"))

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalUserService:     us,
			ExternalLocationAddress: "http://service.url",
		})

		req := httptest.NewRequest("POST", "/someEndpoint", nil)
		req.Header.Set("authorization", bearerToken)
		req.Header.Set("channel-id", channelID)

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		resp := w.Result()

		defer func() { _ = resp.Body.Close() }()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Status code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header")

		expectedJSON := `{"error":"authorization failed: some user service GRPC error"}`
		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})

	t.Run("when user service failed to retrieve Actor and put it in the request context (Basic User not found in repository)", func(t *testing.T) {
		us := new(mocks.ExternalUserServiceMock)
		us.On("ActorFromRequest", bearerToken, ref.ChannelID(channelID), "").
			Return(actor.Actor{}, domain.WrapErrorf(errors.New("record not found"), domain.ErrorCodeUserNotAuthorized, "user could not be authorized"))

		server := NewServer(Config{
			Addr:                    "service.url",
			Logger:                  logger,
			ExternalUserService:     us,
			ExternalLocationAddress: "http://service.url",
		})

		req := httptest.NewRequest("POST", "/someEndpoint", nil)
		req.Header.Set("authorization", bearerToken)
		req.Header.Set("channel-id", channelID)

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

		expectedJSON := `{"error":"user could not be authorized: record not found"}`
		assert.JSONEq(t, expectedJSON, string(b), "response does not match")
	})
}
