package model

import (
	"encoding/json"
	"time"
)

const (
	TYPE_ALIVE = "alive"

	TYPE_AUTH    = "auth"
	TYPE_SIGN_UP = "signup"
	TYPE_SIGN_IN = "signin"

	TYPE_PEER = "peer"
)

// active pool
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// THIS IS INPUT
type PacketWrapper struct {
	Type    string          `json:"type" validate:"required"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// THIS TYPES SEND BY USER
type Packet struct { // SINGLE MESSAGE STORE
	SenderId   string `json:"sender_id" validate:"required"`
	ReceiverId string `json:"receiver_id" validate:"required"`
	Message    string `json:"message" validate:"required"`
}
type ChatHistoryPacket struct { // GET LATEST HISTORY
	SenderId   string    `json:"sender_id" validate:"required"` // THIS IS ALSO THE USER WHO IS REQUESTING FOR MESSAGE
	ReceiverId string    `json:"receiver_id" validate:"required"`
	LastUpdate time.Time `json:"last_update"`
}
type ListPacket struct {
	Uid string `json:"uid" validate:"required"`
}
type PeerPacket struct {
	RemotePeerID string `json:"remote_peer_id"` // "ns" = NOT SPECIFIED || ACTUAL ID
}

// THIS IS THE WRAPPER SEND VIA EVENT
type EventPacket struct {
	Type      string          `json:"type" validate:"required"`
	SessionId string          `json:"session_id" validate:"required"`
	PeerId    string          `json:"peer_id" validate:"required"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

// THESE ARE RESP
type Messages struct { // CHAT HISTORY
	Messages []ChatPacket `json:"messages"`
}
type ChatPacket struct { // SINGLE CHAT BY P2P
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Message    string    `json:"message"`
	Offline    bool      `json:"offline"`
	CreatedAt  time.Time `json:"created_at"`
}
type ActivePool struct { // TOTAL CLIENT ACTIVE ON CURRENT TIME
	Pool []string `json:"pool"`
}
type AuthResponse struct { // AUTH STATUS
	Status bool   `json:"status"`
	Uid    string `json:"uid"`
}
