package server

import "encoding/json"

type response struct {
	Token     string          `json:"token"`
	Payload   json.RawMessage `json:"payload"`
	ErrorCode *string         `json:"error_code,omitempty"`
	Error     *string         `json:"error,omitempty"`
}
