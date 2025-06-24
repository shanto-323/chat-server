package backend

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Manager struct {
	clients ClientList
	mu      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) ServerWS(w http.ResponseWriter, r *http.Request) error {
	socket := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := socket.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	log.Println("New Client: ", conn.RemoteAddr())
	client := NewClient(conn, m)
	m.addClient(client)
	return nil
}

func (m *Manager) addClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[client.id] = client
}

func (m *Manager) removeClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.clients[client.id]; ok {
		client.conn.Close()
		client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "bye bye!"))
		delete(m.clients, client.id)
		return
	}
}
