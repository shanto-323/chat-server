package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"
)

type ClientList map[string]*Client

type Client struct {
	id      string
	conn    *websocket.Conn
	manager *Manager
	msgPool chan IncommingMessage
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	id := ksuid.New().String()
	return &Client{
		id:      id,
		conn:    conn,
		manager: manager,
		msgPool: make(chan IncommingMessage, 1024),
	}
}

func (c *Client) ReadMsg() {
	defer func() {
		c.manager.wg.Done()
		c.manager.removeClient(c)
	}()

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		event := NewEvent(c, payload, c.manager.logger)
		message := event.CreateMessage(payload)
		if message == nil {
			continue
		}
		message.SenderId = c.id

		switch message.MsgType {
		case TYPE_CHAT:
			receiverClient := c.manager.clients[message.ReceiverId]
			event.ChatEvent(receiverClient, message)
		case TYPE_LIST:
			event.ListEvent(message)
		case TYPE_ALIVE:
			waitTime := 30 * time.Second
			c.conn.SetReadDeadline(time.Now().Add(waitTime))
		}

	}
}

func (c *Client) WriteMsg() {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		c.manager.wg.Done()
		c.manager.removeClient(c)
	}()
	for {
		select {
		case msg := <-c.msgPool:
			outgoingMessage := &OutgoingMessage{
				MsgType:  msg.MsgType,
				SenderId: msg.SenderId,
				Payload:  msg.Payload,
			}

			if err := c.conn.WriteJSON(outgoingMessage); err != nil {
				log.Println(err)
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				log.Println(err)
				return
			}
			c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		}
	}
}
