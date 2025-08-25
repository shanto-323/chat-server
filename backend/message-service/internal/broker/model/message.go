package model

import (
	"time"

	"github.com/shanto-323/Chat-Server-1/message-service/internal/database/model"
)

// THIS IS SEND BY CLIENT WITH PACKET
type IncommingMessage struct {
	Message string `json:"message"`
}

// THIS IS INCOMMING WITH PACKET AS PARAM FOR CHAT_HISTORY
type IncommingMessageParam struct {
	LastUpdate time.Time `json:"last_update"`
}

// CLIENT MESSAGE
type Messages struct {
	Messages []model.ChatPacket `json:"messages"`
}
