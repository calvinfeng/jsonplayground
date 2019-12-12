package main

import (
	"encoding/json"
	"fmt"
)

// HTTPErrorResponse is an error JSON response.
type HTTPErrorResponse struct {
	Message string `json:"message,omitempty"`
	Error   error  `json:"error,omitempty"`
}

// HTTPSuccessResponse is a success JSON response.
type HTTPSuccessResponse struct {
	Message     string          `json:"message,omitempty"`
	RequestBody json.RawMessage `json:"request_body,omitempty"`
}

// ValidationError is a map of field to error message.
type ValidationError map[string]string

func (err ValidationError) Error() string {
	var str string
	for field, val := range err {
		str += fmt.Sprintf("%s:%s", field, val)
	}
	return str
}
