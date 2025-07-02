package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"chat_app/backend/logger"
	"chat_app/backend/pkg/storage/redis"

	"github.com/gorilla/websocket"
)

type Manager struct {
	clients     ClientList
	mu          sync.RWMutex
	redisClient *redis.RedisClient
	logger      *logger.ZapLogger
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

func NewManager(ctxP context.Context, redisClient *redis.RedisClient, logger *logger.ZapLogger) *Manager {
	ctx, cancel := context.WithCancel(ctxP)
	return &Manager{
		clients:     make(ClientList),
		redisClient: redisClient,
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
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

	m.logger.Info(fmt.Sprintln("New Client:", conn.RemoteAddr()))

	// add event in client
	client := NewClient(conn, m)
	m.addClient(client)
	m.sendInfo(client)

	m.wg.Add(2)
	go client.ReadMsg()
	go client.WriteMsg()
	return nil
}

func (m *Manager) Shutdown(ctxS context.Context) {
	m.cancel()
	wgDone := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(wgDone)
	}()

	for _, c := range m.clients {
		go func(c *Client) {
			message := IncommingMessage{
				MsgType:    TYPE_CLOSE,
				ReceiverId: c.id,
			}

			select {
			case c.msgPool <- message:
				m.logger.Info("Closing call sended.")
			default:
				m.logger.Error("Buffer is full!")
			}
		}(c)
	}

	select {
	case <-wgDone:
		m.logger.Info("Graceful shutdown")
	case <-ctxS.Done():
		m.logger.Error("Forced shutdown")
	}
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
		m.logger.Info("Info about client.")
	default:
		m.logger.Error("Buffer is full!")
	}
}

func (m *Manager) addClient(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.redisClient.Set(m.ctx, c.conn.LocalAddr().Network(), true)

	m.clients[c.id] = c
	log.Println(m.clients[c.id].id)
}

func (m *Manager) removeClient(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// REMOVE FROM REDIS
	m.redisClient.Remove(m.ctx, c.conn.LocalAddr().Network())

	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "bye bye!"))
	c.conn.Close()
}
