package model

import "encoding/json"

const (
	TYPE_CHAT         = "chat"
	TYPE_CHAT_HISTORY = "chat_history"
	TYPE_LIST         = "list"
	TYPE_ALIVE        = "alive"
)

type Packet struct {
	Type       string          `json:"type" validate:"required"`
	SenderId   string          `json:"sender_id" validate:"required"`
	ReceiverId string          `json:"receiver_id" validate:"required"`
	Payload    json.RawMessage `json:"payload,omitempty"`
}

// THIS WILL SEND VIA EVENT
type EventPacket struct {
	Type      string          `json:"type" validate:"required"`
	SessionId string          `json:"session_id" validate:"required"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}
