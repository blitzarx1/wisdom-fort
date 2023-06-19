package api

import (
	"encoding/json"
)

type Request struct {
	Token   *string         `json:"token"`
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}
