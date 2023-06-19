package api

import "encoding/json"

type Response struct {
	Token     string          `json:"token"`
	Payload   json.RawMessage `json:"payload"`
	ErrorCode *string         `json:"error_code,omitempty"`
	Error     *string         `json:"error,omitempty"`
}
