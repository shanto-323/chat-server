package websocket

import (
	"encoding/json"
	"fmt"

	"chat_app/backend/logger"
)

const (
	TYPE_INFO  = "info"
	TYPE_CHAT  = "chat"
	TYPE_LIST  = "list"
	TYPE_ALIVE = "ping"
	TYPE_CLOSE = "close"
)

type Event interface {
	CreateMessage(payload []byte) (*IncommingMessage, error)
	ChatEvent(c *Client, message IncommingMessage) error
	ListEvent(c *Client, message IncommingMessage) ([]string, error)
	InfoEvent(c *Client) error
	WriteMsg(c *Client, message IncommingMessage) error
}

type event struct {
	logger logger.ZapLogger
}

func NewEvent(logger logger.ZapLogger) Event {
	return &event{
		logger: logger,
	}
}

func (e *event) CreateMessage(payload []byte) (*IncommingMessage, error) {
	var message IncommingMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}
	return &message, nil
}

func (e *event) ChatEvent(c *Client, message IncommingMessage) error {
	select {
	case c.MsgPool <- message:
		return nil
	default:
		return fmt.Errorf("Buffer is full!")
	}
}

func (e *event) ListEvent(c *Client, message IncommingMessage) ([]string, error) {
	var list []string
	clients := c.manager.clients

	for i := range clients {
		if c.manager.clients[i].id == c.id {
			continue
		}
		list = append(list, c.manager.clients[i].id)
	}

	payload, err := json.Marshal(&UserList{IdList: list})
	if err != nil {
		return nil, err
	}
	message.Payload = payload

	select {
	case c.MsgPool <- message:
		return list, nil
	default:
		return nil, fmt.Errorf("Buffer is full!")
	}
}

func (e *event) InfoEvent(c *Client) error {
	message := IncommingMessage{
		MsgType:    TYPE_INFO,
		SenderId:   c.id,
		ReceiverId: c.id,
	}
	payload := UserModel{
		Id:       c.id,
		ConnAddr: c.conn.RemoteAddr().String(),
	}
	payloadJson, _ := json.Marshal(&payload)
	message.Payload = payloadJson

	select {
	case c.MsgPool <- message:
		return nil
	default:
		return fmt.Errorf("Buffer is full!")
	}
}

func (e *event) WriteMsg(c *Client, message IncommingMessage) error {
	outgoingMessage := &OutgoingMessage{
		MsgType:  message.MsgType,
		SenderId: message.SenderId,
		Payload:  message.Payload,
	}

	if err := c.conn.WriteJSON(outgoingMessage); err != nil {
		return err
	}
	return nil
}
