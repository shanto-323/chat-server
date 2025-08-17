package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/client"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/client/model"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/queue"
)

var CLIENT_POOL = map[string]*Client{}

type Manager struct {
	Consumer   queue.Consumer
	UserClient client.UserClient
	mu         *sync.RWMutex
}

func NewManager(ctx context.Context, consumer queue.Consumer) *Manager {
	return &Manager{
		Consumer:   consumer,
		UserClient: client.NewClient(),
		mu:         &sync.RWMutex{},
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
	conn.WriteMessage(websocket.TextMessage, []byte("working"))
	request := model.UserRequest{}
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
		}
		if err := json.Unmarshal(payload, &request); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
		}
		resp, err := m.UserClient.Auth(&request)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
		}
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("success. Client id : %d", resp.ID)))
		break
	}
	go m.addClient(conn)
	// REFACTOR AND REQUEST FOR CACHE
}

func (m *Manager) addClient(conn *websocket.Conn) {
	c := NewClient(conn, m)

	m.mu.Lock()
	CLIENT_POOL[c.ID] = c
	m.mu.Unlock()

	go c.ReadMsg()
	go c.WriteMsg()
}

func (m *Manager) removeClient(c *Client) {
	m.mu.Lock()
	client := CLIENT_POOL[c.ID]
	m.mu.Unlock()

	client.Cancel()

	slog.Info("Client Disconnected !!")
	client.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "bye bye!"))
	client.Conn.Close()
}
