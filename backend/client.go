package backend

import (
	"encoding/json"
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
		msgPool: make(chan IncommingMessage, 1024), // BUFFER SIZE 1024 bytes
	}
}

func (c *Client) ReadMsg() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		var message IncommingMessage
		if err := json.Unmarshal(payload, &message); err != nil {
			log.Println("error marshaling message", err)
			continue
		}
		message.SenderId = c.id

		switch message.MsgType {
		case TYPE_CHAT:
			receiverClient := c.manager.clients[message.ReceiverId]
			select {
			case receiverClient.msgPool <- message:
				log.Println("new message: ", TYPE_CHAT)
			default:
				log.Println("buf size full")
			}
		case TYPE_LIST:
			var list []string
			for i := range c.manager.clients {
				if c.manager.clients[i].id == c.id {
					continue
				}
				list = append(list, c.manager.clients[i].id)
			}

			payload, _ := json.Marshal(&UserList{IdList: list})
			message.Payload = payload

			select {
			case c.msgPool <- message:
				log.Println("new message: ", TYPE_LIST)
			default:
				log.Println("buf size full")
			}
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
