package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/client"
	model "github.com/shanto-323/Chat-Server-1/gateway-1/pkg/client/model"
	model2 "github.com/shanto-323/Chat-Server-1/gateway-1/pkg/connection/model"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/queue"
)

var CLIENT_POOL = map[string]*Client{}

type Manager struct {
	Consumer   queue.Consumer
	UserClient client.UserClient
	mu         *sync.RWMutex
}

func NewManager(ctx context.Context, consumer queue.Consumer) *Manager {
	return &Manager{
		Consumer:   consumer,
		UserClient: client.NewClient(),
		mu:         &sync.RWMutex{},
	}
}

func (m *Manager) ServerWS(w http.ResponseWriter, r *http.Request) {
	socket := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := socket.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if err := m.sendMessage(); err != nil {
		slog.Error(err.Error())
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte("CHAT.APP-1"))
	for {
		if err := m.auth(conn); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
		}
		break
	}

	go m.addClient(conn)
	// REFACTOR AND REQUEST FOR CACHE
}

func (m *Manager) auth(conn *websocket.Conn) error {
	request := model.UserRequest{}
	_, payload, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}
	resp, err := m.UserClient.Auth(&request)
	if err != nil {
		return err
	}
	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("success. Client id : %d", resp.ID)))
	return nil
}

func (m *Manager) addClient(conn *websocket.Conn) {
	c := NewClient(conn, m)

	m.mu.Lock()
	CLIENT_POOL[c.ID] = c
	m.mu.Unlock()

	go c.ReadMsg()
	go c.WriteMsg()
}

func (m *Manager) removeClient(c *Client) {
	m.mu.Lock()
	client := CLIENT_POOL[c.ID]
	m.mu.Unlock()

	client.Cancel()

	slog.Info("Client Disconnected !!")
	client.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "bye bye!"))
	client.Conn.Close()
}

func (m *Manager) sendMessage() error {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	err = m.Consumer.CreateQueue("message.write", true, false)
	if err != nil {
		return err
	}
	err = m.Consumer.CreateQueueBinding("message.write", "gateway.1", "message.service")
	if err != nil {
		return err
	}
	msgChan, err := m.Consumer.Consume("message.write", "", false)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgChan {
			message := model2.Message{}
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				slog.Error(err.Error())
				continue
			}
			c := CLIENT_POOL[message.UserId]

			m.mu.Lock()
			// THIS IS ERROR BECAUSE I DIDNT MAP ACTUAL ID WITH CONN. IT SOPPSE TO STORE IN CLIENT-SERVICE WHICH I DID NOT IMPLEMENTED YET.
			// ALSO NEED A PROPER ERROR HANDLING (i.e CHECK USER IS ALIVE OR NOT)
			c.MsgChan <- &message
			m.mu.Unlock()
		}
	}()
	return nil
}
