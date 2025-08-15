package model

import "time"

type Chat struct {
	ChatId    string    `json:"chat_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type OfflineChat struct {
	SenderId   uint      `json:"sender_id"`
	ReceiverId uint      `json:"receiver_id"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}
