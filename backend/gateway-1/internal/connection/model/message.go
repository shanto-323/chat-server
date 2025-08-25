package model

import "time"

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
