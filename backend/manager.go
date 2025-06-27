package backend

import (
	"encoding/json"
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
	return &Manager{
		clients: make(ClientList),
	}
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
	m.sendInfo(client)

	go client.ReadMsg()
	go client.WriteMsg()
	return nil
}

func (m *Manager) sendInfo(c *Client) {
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
	case c.msgPool <- message:
		log.Println("new message: ", TYPE_INFO)
	default:
		log.Println("buf size full")
	}
}

func (m *Manager) addClient(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[c.id] = c
	log.Println(m.clients[c.id].id)
}

func (m *Manager) removeClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "bye bye!"))
	client.conn.Close()
	delete(m.clients, client.id)
}
