package errx

import (
	"encoding/json"
)

// APIError contains errors about
type APIError struct {
	Err        error `json:"error"`
	StatusCode int   `json:"statusCode"`
}

type APIErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// IsError check if instance is error
func (e *APIError) IsError() bool {
	if e.StatusCode >= 400 {
		return true
	}
	return false
}

func (e *APIError) Error() string {
	return e.Err.Error()
}

// Serialize returns json encoded string
func (e *APIError) Serialize() string {
	if e.IsError() {
		err := APIErrorResponse{
			Message:    e.Err.Error(),
			StatusCode: e.StatusCode,
		}
		// TODO: handle error on serialization error
		content, _ := json.Marshal(err)
		return string(content)
	}
	return "err undefined"
}

// New creates new api error
func New(err error, statusCode int) APIError {
	return APIError{
		Err:        err,
		StatusCode: statusCode,
	}
}
