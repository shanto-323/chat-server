package model

import (
	"time"

	"github.com/gocql/gocql"
)

type Chat struct {
	ChatID         gocql.UUID `json:"chat_id"`
	ConversationID string     `json:"conversation_id"`
	SenderID       string     `json:"sender_id"`
	ReceiverID     string     `json:"receiver_id"`
	Message        string     `json:"message"`
	CreatedAt      time.Time  `json:"created_at"`
	Offline        bool       `json:"offline"`
}
