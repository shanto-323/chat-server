package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"chat_app/backend/logger"
	"chat_app/backend/pkg/storage/redis"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSocket(t *testing.T) {
	logger, _ := logger.NewLogger()
	ctx := context.WithoutCancel(context.Background())

	// New Manager
	redisUrl := "redis://:123456@localhost:6379/0"
	redisMockClient, _ := redis.NewRedisClient(redisUrl, logger)
	m := NewManager(ctx, logger, redisMockClient)

	// Http Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := m.ServerWS(w, r)
		if err != nil {
			assert.Fail(t, "server error")
		}
	}))
	defer server.Close()

	// Connection Upgrade
	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		assert.Fail(t, "dial failed")
	}
	c := NewClient(conn, m)

	t.Run("info", func(t *testing.T) {
		assert.Nil(t, c.manager.event.InfoEvent(c))

		message := IncommingMessage{}
		select {
		case msg := <-c.MsgPool:
			assert.NoError(t, json.Unmarshal(msg.Payload, &message))
		default:
			assert.Fail(t, "cant find message")
		}
		err := c.manager.event.WriteMsg(c, message)
		assert.Nil(t, err)
	})

	t.Run("Chat", func(t *testing.T) {
		m.wg.Add(1)
		message := IncommingMessage{
			MsgType:    TYPE_CHAT,
			ReceiverId: c.id,
			SenderId:   c.id,
		}
		err = c.manager.event.ChatEvent(c, message)
		assert.Nil(t, err)

		err := c.manager.event.WriteMsg(c, message)
		assert.Nil(t, err)
	})
}
