package backend

import (
	"fmt"
	"log"

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
		msgPool: make(chan []byte, 1024), // BUFFER SIZE 1024 bytes
	}
}

func (c *Client) ReadMsg() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		select {
		case c.msgPool <- msg:
			log.Println("new message: ", string(msg))
		default:
			log.Println("buf size full")
		}
	}
}

func (c *Client) ReadMsgForClient(id string, msg string) error {
	client := c.manager.clients[id]
	if client == nil {
		log.Println("user not found")
		return fmt.Errorf("user not found")
	}

	messageBuf := []byte(msg)
	select {
	case client.msgPool <- messageBuf:
		log.Println("incoming message")
		return nil
	default:
		log.Println("something wrong")
		return fmt.Errorf("buf size full")
	}
}

func (c *Client) WriteMsg() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for msg := range c.msgPool {
		log.Println(string(msg))
	}
}
