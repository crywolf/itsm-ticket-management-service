package converters

import (
	"io/ioutil"
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
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

func (c BasePayloadConverter) readBody(r *http.Request) ([]byte, error) {
	defer func() { _ = r.Body.Close() }()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "could not read request body"
		c.logger.Errorw(msg, "error", err)
		err = presenters.WrapErrorf(err, http.StatusInternalServerError, msg)
		return reqBody, err
	}

	return reqBody, nil
}
