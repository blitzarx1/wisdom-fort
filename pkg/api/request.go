package api

import (
	"encoding/json"
)

type Request struct {
	Token   *string         `json:"token"`
	Action  Action          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}
