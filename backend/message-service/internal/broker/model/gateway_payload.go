package model

type GatewayPayload struct {
	SessionId string   `json:"session_id"`
	Data      string   `json:"data"`
	Pool      []string `json:"pool"`
}
