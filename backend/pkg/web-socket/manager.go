package websocket

import (
	"context"
	"net/http"
	"sync"

	"chat_app/backend/logger"
	"chat_app/backend/pkg/storage/redis"

	"github.com/gorilla/websocket"
)

type Manager struct {
	Event  Event
	Logger logger.ZapLogger
	Ctx    context.Context
	Cancel context.CancelFunc
	mu     sync.RWMutex
	wg     sync.WaitGroup
}

func NewManager(ctx context.Context, l logger.ZapLogger, r redis.RedisClient) *Manager {
	ctx, cancel := context.WithCancel(ctx)
	event := NewEvent(l, r, nil)
	return &Manager{
		Logger: l,
		Ctx:    ctx,
		Cancel: cancel,
		Event:  event,
		wg:     sync.WaitGroup{},
	}
}

func (m *Manager) ServerWS(w http.ResponseWriter, r *http.Request) {
	socket := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := socket.Upgrade(w, r, nil)
	if err != nil {
		m.Logger.Error(err.Error())
	}

	m.addClient(conn)
}

func (m *Manager) addClient(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	c := NewClient(conn, m)

	m.wg.Add(2)
	go c.ReadMsg()
	go c.WriteMsg()
	m.Event.AddClient(c)
}

func (m *Manager) Shutdown(ctxS context.Context) {
	l := m.Logger

	defer m.Cancel()
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		l.Info("Graceful shutdown")
	case <-ctxS.Done():
		l.Error("Forced shutdown")
	}
}

func (m *Manager) removeClient(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "bye bye!"))
	c.conn.Close()
}
