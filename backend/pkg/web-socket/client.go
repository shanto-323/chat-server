package websocket

import (
	"time"

	"chat_app/backend/pkg/model"

	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"
)

type Client struct {
	id      string
	conn    *websocket.Conn
	MsgPool chan model.MessagePacket
	manager *Manager
}

func NewClient(conn *websocket.Conn, m *Manager) *Client {
	id := ksuid.New().String()
	return &Client{
		id:      id,
		conn:    conn,
		MsgPool: make(chan model.MessagePacket, 1024),
		manager: m,
	}
}

func (c *Client) ReadMsg() {
	defer func() {
		c.manager.wg.Done()
		c.manager.removeClient(c)
	}()
	m := c.manager

	for {
		select {
		case <-m.Ctx.Done():
			return
		default:
			_, payload, err := c.conn.ReadMessage()
			if err != nil {
				return
			}

			message, err := m.Event.CreateMessage(payload)
			if err != nil {
				m.Logger.Error(err.Error())
				continue
			}

			switch message.MsgType {
			case TYPE_CHAT:
				if err := m.Event.ChatEvent(message); err != nil {
					m.Logger.Error(err.Error())
				}
			case TYPE_LIST:
				m.Logger.Info("List!!")
				if err := m.Event.ListEvent(c, message); err != nil {
					m.Logger.Error(err.Error())
				}
			case TYPE_ALIVE:
				waitTime := 30 * time.Second
				c.conn.SetReadDeadline(time.Now().Add(waitTime))
			}
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
	m := c.manager

	for {
		select {
		case <-c.manager.Ctx.Done():
			return
		case msg := <-c.MsgPool:
			m.Event.WriteMsg(c, msg)
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				m.Logger.Error(err.Error())
				return
			}
			c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		}
	}
}
