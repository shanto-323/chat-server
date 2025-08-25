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
	data "github.com/shanto-323/Chat-Server-1/gateway-1/data/remote"
	dataModel "github.com/shanto-323/Chat-Server-1/gateway-1/data/remote/model"
	"github.com/shanto-323/Chat-Server-1/gateway-1/internal/broker"
	"github.com/shanto-323/Chat-Server-1/gateway-1/internal/connection/model"
)

type Manager struct {
	Publisher   broker.Publisher
	Consumer    broker.Consumer
	UserClient  data.UserClient
	CacheClient data.CacheClient
	ClientPool  map[string]*Client
	mu          *sync.RWMutex
}

func NewManager(ctx context.Context, p broker.Publisher, c broker.Consumer) *Manager {
	return &Manager{
		Publisher:   p,
		Consumer:    c,
		UserClient:  data.NewClient(),
		CacheClient: data.NewCacheClient(),
		ClientPool:  map[string]*Client{},
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
	request := dataModel.UserRequest{}
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

func (m *Manager) addCache(conn *websocket.Conn, user *dataModel.User) error {
	userId := strconv.FormatInt(int64(user.ID), 10)
	c := NewClient(conn, m, userId)

	req := dataModel.ConnRequest{
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
	m.ClientPool[c.SessionId] = c
	m.mu.Unlock()

	go c.ReadMsg()
	go c.WriteMsg()
}

func (m *Manager) removeClient(c *Client) {
	m.mu.Lock()
	client := m.ClientPool[c.SessionId]
	m.mu.Unlock()

	slog.Info("REMOVE CONN", "ID", c.ClientId, "SESSION_ID", c.SessionId)
	req := dataModel.ConnRequest{
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
	consumer, err := m.Consumer.Consume()
	if err != nil {
		return err
	}

	go func() {
		for msg := range consumer {
			packet := model.EventPacket{}
			if err := json.Unmarshal(msg.Body, &packet); err != nil {
				slog.Error(err.Error())
				continue
			}

			c, exists := m.ClientPool[packet.SessionId]
			if !exists {
				slog.Error("MANAGER", "User not exist", err.Error())
			}

			m.mu.Lock()
			c.MsgChan <- &packet
			m.mu.Unlock()
		}
	}()
	return nil
}
