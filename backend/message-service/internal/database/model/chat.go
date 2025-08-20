package model

import "time"

type Chat struct {
	ChatID     string    `json:"chat_id"`
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}
