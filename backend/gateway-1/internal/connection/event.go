package connection

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"github.com/shanto-323/Chat-Server-1/gateway-1/internal/connection/model"
)

type Event interface {
	SendPayload(conn *websocket.Conn, messageType string, payload json.RawMessage) error
	AddCache(conn *websocket.Conn, m *Manager, client *Client) error
	EventProcess(d amqp091.Delivery, m *Manager) error
}

type event struct{}

func NewEvent() Event {
	return &event{}
}

func (e *event) SendPayload(conn *websocket.Conn, messageType string, payload json.RawMessage) error {
	message := model.PacketWrapper{
		Type:    messageType,
		Payload: payload,
	}

	if err := conn.WriteJSON(message); err != nil {
		return err
	}
	return nil
}

func (e *event) AddCache(conn *websocket.Conn, m *Manager, client *Client) error {
	if err := m.CacheClient.AddActiveUser(client.ClientId, client.SessionId); err != nil {
		return err
	}

	authResponse := model.AuthResponse{
		Status: true,
		Uid:    client.ClientId,
	}
	payload, err := json.Marshal(&authResponse)
	if err != nil {
		return err
	}

	return m.event.SendPayload(conn, model.TYPE_AUTH, payload)
}

func (e *event) EventProcess(d amqp091.Delivery, m *Manager) error {
	packet := model.EventPacket{}
	if err := json.Unmarshal(d.Body, &packet); err != nil {
		return err
	}

	c, exists := m.ClientPool[packet.SessionId]
	if !exists {
		return fmt.Errorf("user does not exist%s", packet.SessionId)
	}

	m.mu.Lock()
	c.MsgChan <- &packet
	m.mu.Unlock()

	return nil
}
