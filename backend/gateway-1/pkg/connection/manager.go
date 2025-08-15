package connection

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/model"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/queue"
)

type Manager struct {
	// Event       Event
	Ctx      context.Context
	Cancel   context.CancelFunc
	mu       sync.RWMutex
	wg       sync.WaitGroup
	Consumer queue.Consumer
}

func NewManager(ctx context.Context, consumer queue.Consumer) *Manager {
	ctx, cancel := context.WithCancel(ctx)
	return &Manager{
		Ctx:      ctx,
		Cancel:   cancel,
		wg:       sync.WaitGroup{},
		Consumer: consumer,
	}
}

func (m *Manager) ServerWS(w http.ResponseWriter, r *http.Request) {
	socket := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := socket.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err := m.Consumer.SendMessage(context.Background(), "", "", amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Persistent,
		Body:         []byte(`hello new socket`),
	}); err != nil {
		slog.Error(err.Error())
	}

	slog.Info("Message sent ...")

	action := model.UserAction{}
	if err := conn.ReadJSON(&action); err != nil {
		slog.Error("invalid action: " + err.Error())
		conn.Close()
		return
	}

	// m.addClient(conn, resp)
}

// func (m *Manager) addClient(conn *websocket.Conn, resp *model.User) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	c := NewClient(conn, m)
// 	c.id = resp.ID

// 	m.wg.Add(2)
// 	go c.ReadMsg()
// 	go c.WriteMsg()
// 	m.Event.AddClient(c)
// }

func (m *Manager) Shutdown(ctxS context.Context) {
	defer m.Cancel()
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("Graceful shutdown")
	case <-ctxS.Done():
		slog.Error("Forced shutdown")
	}
}

func (m *Manager) removeClient(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// if err := m.Event.RemoveClient(c); err != nil {
	// 	slog.Error(err.Error())
	// }

	slog.Info("Client Disconnected !!")
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "bye bye!"))
	c.conn.Close()
}
