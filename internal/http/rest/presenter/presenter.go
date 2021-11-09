package presenter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// Presenter provides REST responses
type Presenter interface {
	// WriteError replies to the request with the specified error message and HTTP code.
	// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
	// The error message should be plain text.
	WriteError(w http.ResponseWriter, error string, code int)

	// WriteJSON encodes 'v' to JSON and writes it to the 'w'. Also sets correct Content-Type header.
	// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
	WriteJSON(w http.ResponseWriter, v interface{})
}

// NewPresenter creates a presentation service
func NewPresenter(logger *zap.SugaredLogger, serverAddr string) Presenter {
	return &presenter{
		logger:     logger,
		serverAddr: serverAddr,
	}
}

type presenter struct {
	logger     *zap.SugaredLogger
	serverAddr string
}

func (p presenter) WriteError(w http.ResponseWriter, error string, code int) {
	p.sendErrorJSON(w, error, code)
}

func (p presenter) WriteJSON(w http.ResponseWriter, v interface{}) {
	p.encodeJSON(w, v)
}

// sendErrorJSON replies to the request with the specified error message and HTTP code.
// It encodes error string as JSON object {"error":"error_string"} and sets correct header.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
// The error message should be plain text.
func (p presenter) sendErrorJSON(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errorJSON, _ := json.Marshal(error)
	_, _ = fmt.Fprintf(w, `{"error":%s}`+"\n", errorJSON)
}

// encodeJSON encodes 'v' to JSON and writes it to the 'w'. Also sets correct Content-Type header.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
func (p presenter) encodeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		eMsg := "could not encode JSON response"
		p.logger.Errorw(eMsg, "error", err)
		p.WriteError(w, eMsg, http.StatusInternalServerError)
		return
	}
}
