package model

import (
	"encoding/json"
	"time"
)

const (
	TYPE_CHAT         = "chat"
	TYPE_CHAT_HISTORY = "chat_history"
	TYPE_LIST         = "list"
	TYPE_ALIVE        = "alive"

	TYPE_AUTH    = "auth"
	TYPE_SIGN_UP = "signup"
	TYPE_SIGN_IN = "signin"
)

// THIS IS INPUT && OUTPUT MODEL
type PacketWrapper struct {
	Type    string          `json:"type" validate:"required"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// PAYLOAD TYPES --->

// pub/sub model
type EventPacket struct {
	Type      string          `json:"type" validate:"required"`
	SessionId string          `json:"session_id" validate:"required"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

// active pool
type ActivePool struct {
	Pool []struct {
		Uid      string `json:"uid"`
		Username string `json:"username"`
	} `json:"pool"`
}

// auth status
type AuthResponse struct {
	Status bool `json:"status"`
}

// messages
type Messages struct {
	Messages []ChatPacket `json:"messages"`
}

type ChatPacket struct {
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Message    string    `json:"message"`
	Offline    bool      `json:"offline"`
	CreatedAt  time.Time `json:"created_at"`
}
