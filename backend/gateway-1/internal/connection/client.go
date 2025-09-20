package connection

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
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

func NewClient(conn *websocket.Conn, m *Manager, uid uint) *Client {
	clientId := strconv.FormatInt(int64(uid), 10)
	sessionId := ksuid.New().String()

	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ClientId:  clientId,
		SessionId: sessionId,
		Conn:      conn,
		Manager:   m,
		Cancel:    cancel,
		Ctx:       ctx,
		MsgChan:   make(chan *model.EventPacket),
	}
}

func (c *Client) ReadMsg() {
	defer func() {
		c.Manager.removeClient(c)
	}()

	m := c.Manager
	conn := c.Conn
	for {
		select {
		case <-c.Ctx.Done():
			return
		default:
			_, payload, err := conn.ReadMessage()
			if err != nil {
				return
			}

			packet := model.PacketWrapper{}
			if err := json.Unmarshal(payload, &packet); err != nil {
				continue
			}
			slog.Info("NEW PACKET", c.ClientId, packet)

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

func (c *Client) WriteMsg() {
	defer func() {
		c.Manager.removeClient(c)
	}()
	ticker := time.NewTicker(10 * time.Second)
	conn := c.Conn

	for {
		select {
		case <-c.Ctx.Done():
			return
		case _ = <-c.MsgChan:
		//Message need to send
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				return
			}
		default:
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		}
	}
}
