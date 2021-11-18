package converters

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
	"github.com/go-openapi/runtime/middleware/header"
	"go.uber.org/zap"
)

// NewBasePayloadConverter returns input payload converting service with basic functionality
func NewBasePayloadConverter(logger *zap.SugaredLogger) *BasePayloadConverter {
	return &BasePayloadConverter{
		logger: logger,
	}
}

// BasePayloadConverter must be included in all derived converters via object composition
type BasePayloadConverter struct {
	logger *zap.SugaredLogger
}

func (c BasePayloadConverter) unmarshalFromBody(r *http.Request, dst interface{}) error {
	defer func() { _ = r.Body.Close() }()

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			return presenters.NewErrorf(http.StatusUnsupportedMediaType, "Content-Type header is not application/json")
		}
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return presenters.WrapErrorf(err, http.StatusBadRequest, "Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error for syntax errors in the JSON.
		// There is an open issue regarding this at https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return presenters.WrapErrorf(err, http.StatusBadRequest, "Request body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			return presenters.NewErrorf(http.StatusBadRequest, "Request body contains an invalid value for the %q field (type: %s, value: %s)", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value)

		// Catch the error caused by extra unexpected fields in the request body.
		// There is an open issue at https://github.com/golang/go/issues/29035 regarding turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return presenters.WrapErrorf(err, http.StatusBadRequest, "Request body contains unknown field %s", fieldName)

		case errors.Is(err, io.EOF):
			return presenters.WrapErrorf(err, http.StatusBadRequest, "Request body must not be empty")

		default:
			return err
		}
	}

	// Call decode again, if the request body only contained a single JSON object this will return an io.EOF error.
	// So if we get anything else, we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return presenters.NewErrorf(http.StatusBadRequest, "Request body must only contain a single JSON object")
	}

	return nil
}
