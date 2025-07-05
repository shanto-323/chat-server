package websocket

import (
	"context"
	"net/http"
	"sync"
	"time"

	"chat_app/backend/logger"
	"chat_app/backend/pkg/storage/redis"

	"github.com/gorilla/websocket"
)

type ClientList map[string]*Client

type Manager struct {
	redisClient redis.RedisClient
	clients     ClientList // Swap it with database
	event       Event
	logger      logger.ZapLogger
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
	wg          sync.WaitGroup
}

func NewManager(ctx context.Context, l logger.ZapLogger, r redis.RedisClient) *Manager {
	ctx, cancel := context.WithCancel(ctx)
	return &Manager{
		redisClient: r,
		logger:      l,
		ctx:         ctx,
		cancel:      cancel,
		event:       NewEvent(l),
		clients:     make(ClientList),
		wg:          sync.WaitGroup{},
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

	return m.addClient(conn)
}

func (m *Manager) addClient(conn *websocket.Conn) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	c := NewClient(conn, m)

	m.wg.Add(2)
	go c.ReadMsg()
	go c.WriteMsg()
	m.redisClient.Set(m.ctx, c.conn.LocalAddr().Network(), true)
	m.redisClient.SetList(m.ctx, redis.LIST_KEY, c.id)
	return m.event.InfoEvent(c)
}

func (m *Manager) Shutdown(ctxS context.Context) {
	m.cancel()
	wgDone := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(wgDone)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clients, err := m.redisClient.GetList(ctx, redis.LIST_KEY)
	if err != nil {
		m.logger.Error(err.Error())
	}

	for _, cid := range clients {
		go func(id string) {
			message := IncommingMessage{
				MsgType:    TYPE_CLOSE,
				ReceiverId: id,
			}
			c := m.clients[cid]
			m.event.WriteMsg(c, message)
		}(cid)
	}

	select {
	case <-wgDone:
		m.redisClient.Close()
		m.logger.Info("Graceful shutdown")
	case <-ctxS.Done():
		m.logger.Error("Forced shutdown")
	}
}

func (m *Manager) removeClient(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx := context.WithoutCancel(context.Background())
	m.redisClient.Remove(ctx, c.conn.LocalAddr().Network())

	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "bye bye!"))
	c.conn.Close()
}
