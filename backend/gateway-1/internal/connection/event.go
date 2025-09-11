package connection

import (
	"encoding/json"
	"log/slog"

	"github.com/shanto-323/Chat-Server-1/gateway-1/internal/connection/model"
)

type Event interface {
	ChatEvent(payload json.RawMessage)
}

type event struct {
	Client *Client
}

func NewEvent(c *Client) Event {
	return &event{
		Client: c,
	}
}

func (e *event) ChatEvent(payload json.RawMessage) {
	message := model.ChatPacket{}
	if err := json.Unmarshal(payload, &message); err != nil {
		slog.Error(err.Error())
		return
	}

	if err := e.Client.Conn.WriteJSON(message); err != nil {
		slog.Error(err.Error())
		return
	}
}
