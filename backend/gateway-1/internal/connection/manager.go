package connection

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	data "github.com/shanto-323/Chat-Server-1/gateway-1/data/remote"
	"github.com/shanto-323/Chat-Server-1/gateway-1/internal/broker"
	"github.com/shanto-323/Chat-Server-1/gateway-1/internal/connection/model"
)

type Manager struct {
	Publisher   broker.Publisher
	Consumer    broker.Consumer
	UserClient  data.UserClient
	CacheClient data.CacheClient
	event       Event

	mu         *sync.RWMutex
	ClientPool map[string]*Client
}

func NewManager(ctx context.Context, p broker.Publisher, c broker.Consumer) *Manager {
	return &Manager{
		Publisher:   p,
		Consumer:    c,
		UserClient:  data.NewClient(),
		CacheClient: data.NewCacheClient(),
		event:       NewEvent(),

		mu:         &sync.RWMutex{},
		ClientPool: map[string]*Client{},
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

	slog.Info("NEW CONN", "ip-port", conn.RemoteAddr())
	for {
		if err := m.auth(conn); err != nil {
			authResponse := model.AuthResponse{
				Status: false,
			}
			payload, err := json.Marshal(&authResponse)
			if err != nil {
				slog.Error(err.Error())
			}

			if err := m.event.SendPayload(conn, model.TYPE_AUTH, payload); err != nil {
				slog.Error(err.Error())
			}
			continue
		}
		break
	}
}

func (m *Manager) auth(conn *websocket.Conn) error {
	_, payload, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	resp, err := m.UserClient.Auth(payload)
	if err != nil {
		return err
	}

	client := NewClient(conn, m, resp.ID)
	if err := m.event.AddCache(conn, m, client); err != nil {
		return err
	}

	go m.addClient(client)
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
	if err := m.CacheClient.RemoveActiveUser(c.ClientId, c.SessionId); err != nil {
		slog.Error(err.Error())
	}

	client.Cancel()

	slog.Info("Client Disconnected !!")
	if err := client.Conn.Close(); err != nil {
		slog.Error(err.Error())
	}
}

func (m *Manager) ConsumerStream() error {
	consumer, err := m.Consumer.Consume()
	if err != nil {
		return err
	}

	go func() {
		for d := range consumer {
			go func() {
				if err := m.event.EventProcess(d, m); err != nil {
					slog.Error(err.Error())
				}
			}()
		}
	}()
	return nil
}
