package connection

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/ksuid"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/connection/model"
)

type Client struct {
	ClientId  string
	SessionId string
	Conn      *websocket.Conn
	Manager   *Manager
	Cancel    context.CancelFunc
	Ctx       context.Context
	MsgChan   chan *model.Message
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
		MsgChan:   make(chan *model.Message),
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
					conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
					continue
				}

				packet := model.Packet{}
				if err := json.Unmarshal(payload, &packet); err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
					continue
				}

				switch packet.Type {
				case model.TYPE_CHAT:
					m.Consumer.SendMessage(context.Background(), "message.service", "incomming.message", amqp091.Publishing{
						ContentType:  "text/plain",
						DeliveryMode: amqp091.Persistent,
						Body:         payload,
					})
					conn.WriteMessage(websocket.TextMessage, []byte(packet.Data))
				case model.TYPE_LIST:
					for cp := range CLIENT_POOL {
						conn.WriteMessage(websocket.TextMessage, []byte(cp))
					}

				case model.TYPE_ALIVE:
					waitTime := 30 * time.Second
					conn.SetReadDeadline(time.Now().Add(waitTime))
				}
			}
		}
	}
}

func (c *Client) WriteMsg() {
	ticker := time.NewTicker(10 * time.Second)
	// m := c.Manager
	conn := c.Conn

	for {
		select {
		case <-c.Ctx.Done():
			return
		case msg := <-c.MsgChan:
			conn.WriteMessage(websocket.TextMessage, []byte("incomming --"+msg.Messsage))
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				return
			}
		default:
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		}
	}
}
