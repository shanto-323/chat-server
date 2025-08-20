package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	model2 "github.com/shanto-323/Chat-Server-1/gateway-1/pkg/connection/model"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/queue"
	client "github.com/shanto-323/Chat-Server-1/gateway-1/pkg/remote"
	model "github.com/shanto-323/Chat-Server-1/gateway-1/pkg/remote/model"
)

var CLIENT_POOL = map[string]*Client{}

type Manager struct {
	Consumer    queue.Consumer
	UserClient  client.UserClient
	CacheClient client.CacheClient
	mu          *sync.RWMutex
}

func NewManager(ctx context.Context, consumer queue.Consumer) *Manager {
	return &Manager{
		Consumer:    consumer,
		UserClient:  client.NewClient(),
		CacheClient: client.NewCacheClient(),
		mu:          &sync.RWMutex{},
	}
}

func (m *Manager) ServerWS(w http.ResponseWriter, r *http.Request) {
	socket := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := socket.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if err := m.sendMessage(); err != nil {
		slog.Error(err.Error())
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte("CHAT.APP-1"))
	for {
		if err := m.auth(conn); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
		}
		break
	}
}

func (m *Manager) auth(conn *websocket.Conn) error {
	request := model.UserRequest{}
	_, payload, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}
	resp, err := m.UserClient.Auth(&request)
	if err != nil {
		return err
	}
	if err := m.addCache(conn, resp); err != nil {
		return err
	}

	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("success. Client id : %d", resp.ID)))
	return nil
}

func (m *Manager) addCache(conn *websocket.Conn, user *model.User) error {
	userId := strconv.FormatInt(int64(user.ID), 10)
	c := NewClient(conn, m, userId)

	req := model.ConnRequest{
		ID:        c.ClientId,
		SessionId: c.SessionId,
		GatewayId: "gateway.1",
	}

	if err := m.CacheClient.AddActiveUser(&req); err != nil {
		return err
	}

	go m.addClient(c)
	return nil
}

func (m *Manager) addClient(c *Client) {
	slog.Info("NEW CLIENT", "ID", c.ClientId, "SESSION_ID", c.SessionId)
	m.mu.Lock()
	CLIENT_POOL[c.SessionId] = c
	m.mu.Unlock()

	go c.ReadMsg()
	go c.WriteMsg()
}

func (m *Manager) removeClient(c *Client) {
	m.mu.Lock()
	client := CLIENT_POOL[c.SessionId]
	m.mu.Unlock()

	slog.Info("REMOVE CONN", "ID", c.ClientId, "SESSION_ID", c.SessionId)
	req := model.ConnRequest{
		ID:        c.ClientId,
		SessionId: c.SessionId,
		GatewayId: "gateway.1",
	}

	if err := m.CacheClient.RemoveActiveUser(&req); err != nil {
		slog.Error(err.Error())
	}

	client.Cancel()

	slog.Info("Client Disconnected !!")
	client.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "bye bye!"))
	client.Conn.Close()
}

func (m *Manager) sendMessage() error {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	err = m.Consumer.CreateQueue("message.write", true, false)
	if err != nil {
		return err
	}
	err = m.Consumer.CreateQueueBinding("message.write", "gateway.1", "message.service")
	if err != nil {
		return err
	}
	msgChan, err := m.Consumer.Consume("message.write", "", false)
	if err != nil {
		return err
	}

	go func() {
		slog.Info("SEND MESSAGE RUNNING")
		for msg := range msgChan {
			packet := model2.ConsumePacket{}
			if err := json.Unmarshal(msg.Body, &packet); err != nil {
				slog.Error(err.Error())
				continue
			}
			c := CLIENT_POOL[packet.SessionId]

			m.mu.Lock()
			c.MsgChan <- packet.Data
			m.mu.Unlock()
		}
	}()
	return nil
}
