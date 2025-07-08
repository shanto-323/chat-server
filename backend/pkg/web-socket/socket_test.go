package websocket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"chat_app/backend/logger"
	"chat_app/backend/pkg/model"
	"chat_app/backend/pkg/storage/redis"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSocket(t *testing.T) {
	// logger
	logger, err := logger.NewLogger()
	if err != nil {
		assert.Fail(t, err.Error())
	}
	defer logger.Sync()

	// ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Redis
	redisUrl := "redis://:123456@localhost:6379/0"
	redisMockClient, _ := redis.NewRedisClient(redisUrl)

	// New Manager
	m := NewManager(ctx, logger, redisMockClient)

	// Http Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.ServerWS(w, r)
	}))
	defer server.Close()

	// Connection Upgrade
	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		assert.Fail(t, "dial failed")
	}

	// New Client
	c := NewClient(conn, m)

	t.Run("info", func(t *testing.T) {
		m.Event.AddClient(c)
		select {
		case msg := <-c.MsgPool:
			assert.Equal(t, TYPE_INFO, msg.MsgType)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "cant find message")
		}
	})

	t.Run("chat", func(t *testing.T) {
		message := model.MessagePacket{
			MsgType:    TYPE_CHAT,
			ReceiverId: c.id,
			Payload:    nil,
		}

		err = m.Event.ChatEvent(&message)
		assert.Nil(t, err)

		select {
		case msg := <-c.MsgPool:
			assert.Equal(t, TYPE_CHAT, msg.MsgType)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "cant find message")
		}
	})

	t.Run("list", func(t *testing.T) {
		message := model.MessagePacket{
			MsgType:    TYPE_LIST,
			ReceiverId: c.id,
			Payload:    nil,
		}

		err = m.Event.ListEvent(c, &message)
		assert.Nil(t, err)
		select {
		case msg := <-c.MsgPool:
			assert.Equal(t, TYPE_LIST, msg.MsgType)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "cant find message")
		}
	})
}
