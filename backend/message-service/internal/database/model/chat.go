package model

import "time"

type Chat struct {
	SenderId   uint      `json:"sender_id"`
	ReceiverId uint      `json:"receiver_id"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}
