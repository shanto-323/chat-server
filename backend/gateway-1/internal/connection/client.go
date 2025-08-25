package connection

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"
	"github.com/shanto-323/Chat-Server-1/gateway-1/internal/connection/model"
)

type Client struct {
	ClientId  string
	SessionId string
	Conn      *websocket.Conn
	Manager   *Manager
	Cancel    context.CancelFunc
	Ctx       context.Context
	MsgChan   chan *model.EventPacket
}

func NewClient(conn *websocket.Conn, m *Manager, clientId string) *Client {
	id := ksuid.New().String()
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ClientId:  clientId,
		SessionId: id,
		Conn:      conn,
		Manager:   m,
		Cancel:    cancel,
		Ctx:       ctx,
		MsgChan:   make(chan *model.EventPacket, 1024),
	}
}

func (c *Client) ReadMsg() {
	m := c.Manager
	conn := c.Conn
	for {
		select {
		case <-c.Ctx.Done():
			return
		default:
			{
				_, payload, err := conn.ReadMessage()
				if err != nil {
					m.removeClient(c)
					continue
				}

				packet := model.Packet{}
				if err := json.Unmarshal(payload, &packet); err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
					continue
				}

				if packet.Type == model.TYPE_ALIVE {
					waitTime := 30 * time.Second
					conn.SetReadDeadline(time.Now().Add(waitTime))
				} else {
					if err := m.Publisher.Publish(payload); err != nil {
						slog.Error("CLIENT", "err", err.Error())
						continue
					}
				}

			}
		}
	}
}

func (c *Client) WriteMsg() {
	ticker := time.NewTicker(10 * time.Second)
	conn := c.Conn

	for {
		select {
		case <-c.Ctx.Done():
			return
		case msg := <-c.MsgChan:
			slog.Info("INFO", "msg", msg)
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				return
			}
		default:
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		}
	}
}
