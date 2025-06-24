package backend

import (
	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"
)

type ClientList map[string]*Client

type Client struct {
	id      string
	conn    *websocket.Conn
	manager *Manager
	msgPool chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	id := ksuid.New().String()
	return &Client{
		id:      id,
		conn:    conn,
		manager: manager,
		msgPool: make(chan []byte),
	}
}
