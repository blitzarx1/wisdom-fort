package client

import (
	"encoding/json"
)

type request struct {
	Token   *string         `json:"token"`
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}
