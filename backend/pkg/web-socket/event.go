package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"chat_app/backend/logger"
	"chat_app/backend/pkg/model"
	"chat_app/backend/pkg/storage/database"
	"chat_app/backend/pkg/storage/redis"
)

const (
	TYPE_INFO  = "info"
	TYPE_CHAT  = "chat"
	TYPE_LIST  = "list"
	TYPE_ALIVE = "ping"
	TYPE_CLOSE = "close"
)

type Event interface {
	CreateMessage(payload []byte) (*model.MessagePacket, error)
	AddClient(c *Client)
	ChatEvent(msg *model.MessagePacket) error
	ListEvent(c *Client, message *model.MessagePacket) error
	WriteMsg(c *Client, message model.MessagePacket) error
	CloseEvent(m *Manager, ctx context.Context) error
}

type event struct {
	logger      logger.ZapLogger
	redisClient redis.RedisClient
	repository  database.Repository
	cPool       map[string]*Client
}

func NewEvent(l logger.ZapLogger, r redis.RedisClient, repo database.Repository) Event {
	return &event{
		logger:      l,
		redisClient: r,
		repository:  repo,
		cPool:       make(map[string]*Client),
	}
}

func (e *event) AddClient(c *Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	e.redisClient.SAdd(ctx, c.id, true)
	e.cPool[c.id] = c

	message := model.MessagePacket{
		MsgType:    TYPE_INFO,
		SenderId:   c.id,
		ReceiverId: c.id,
		Payload:    json.RawMessage{},
	}

	payload := model.Client{
		ID: c.id,
	}
	payloadJson, _ := json.Marshal(&payload)
	message.Payload = payloadJson

	select {
	case c.MsgPool <- message:
		e.logger.Info("New Connection....")
	default:
		e.logger.Error("Buffer is full")
	}
}

func (e *event) CreateMessage(payload []byte) (*model.MessagePacket, error) {
	var message model.MessagePacket
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (e *event) ChatEvent(msg *model.MessagePacket) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	alive, err := e.redisClient.IsMember(ctx, msg.ReceiverId)
	if err != nil {
		e.logger.Error(err.Error())
		return err
	}

	if !alive {
		// ScyllaDb For Message Queue
		return nil
	}

	client, ok := e.cPool[msg.ReceiverId]
	if !ok {
		e.logger.Error("redis logic not aligning!!")
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := e.redisClient.SRem(ctx, msg.ReceiverId); err != nil {
				e.logger.Error(err.Error())
			}
		}()
		return nil
	}

	select {
	case client.MsgPool <- *msg:
		return nil
	default:
		return fmt.Errorf("Buffer is full!")
	}
}

func (e *event) ListEvent(c *Client, message *model.MessagePacket) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clients, err := e.redisClient.SMembers(ctx)
	if err != nil {
		return err
	}

	activePool := model.ActivePool{
		AliveList: clients,
	}

	payload, err := json.Marshal(&activePool)
	if err != nil {
		return err
	}
	message.Payload = payload

	select {
	case c.MsgPool <- *message:
		return nil
	default:
		return fmt.Errorf("Buffer is full!")
	}
}

func (e *event) InfoEvent(c *Client) error {
	message := model.MessagePacket{
		MsgType:    TYPE_INFO,
		SenderId:   c.id,
		ReceiverId: c.id,
	}
	payload := c.id
	payloadJson, _ := json.Marshal(&payload)
	message.Payload = payloadJson

	select {
	case c.MsgPool <- message:
		return nil
	default:
		return fmt.Errorf("Buffer is full!")
	}
}

func (e *event) CloseEvent(m *Manager, ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clients, err := e.redisClient.SMembers(ctx)
	if err != nil {
		return err
	}

	for _, id := range clients {
		client, ok := e.cPool[id]
		if !ok {
			e.logger.Error("redis logic not aligning!!")
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				if err := e.redisClient.SRem(ctx, id); err != nil {
					e.logger.Error(err.Error())
				}
			}()
			return nil
		}

		m.wg.Add(1)
		go func(c *Client) {
			defer m.wg.Done()
			message := model.MessagePacket{
				MsgType:    TYPE_CLOSE,
				ReceiverId: c.id,
			}
			e.WriteMsg(c, message)
		}(client)
	}
	return nil
}

func (e *event) WriteMsg(c *Client, message model.MessagePacket) error {
	if err := c.conn.WriteJSON(message); err != nil {
		return err
	}
	return nil
}
