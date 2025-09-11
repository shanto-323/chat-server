package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gocql/gocql"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/database/model"
)

type MessageRepository interface {
	Close()
	InsertMessage(ctx context.Context, chat *model.Chat) error
	GetMessageFromBucket(ctx context.Context, conversation_id string, createdAt time.Time) ([]*model.ChatPacket, error)
	GetLatestMessageFromBucket(ctx context.Context, conversation_id string) (*model.ChatPacket, error)
}

type scyllaRepository struct {
	session *gocql.Session
}

func NewUserRepository(url string) (MessageRepository, error) {
	cluster := gocql.NewCluster(url)
	cluster.Port = 9042
	cluster.Keyspace = "cluster"
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &scyllaRepository{
		session: session,
	}, nil
}

func (s *scyllaRepository) Close() {
	s.session.Close()
}

func (s *scyllaRepository) InsertMessage(ctx context.Context, chat *model.Chat) error {
	query := `
        INSERT INTO chat_history(
            chat_id, conversation_id, payload , created_at
        ) VALUES (?, ?, ?, ?)
    `

	// JSON Payload CREATION
	payload := chat.Payload
	blob, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return s.session.Query(query,
		chat.ChatID,
		chat.ConversationID,
		blob,
		chat.Payload.CreatedAt,
	).WithContext(ctx).Exec()
}

func (s *scyllaRepository) GetMessageFromBucket(ctx context.Context, conversation_id string, createdAt time.Time) ([]*model.ChatPacket, error) {
	query := `
		SELECT payload
		FROM chat_history
		WHERE conversation_id = ?
		  AND created_at < ?
		ORDER BY created_at DESC
		LIMIT 10
	`
	iter := s.session.Query(query, conversation_id, createdAt).WithContext(ctx).Iter()

	var chatHistory []*model.ChatPacket
	var payloadBytes []byte

	for iter.Scan(&payloadBytes) {
		chatPacket := model.ChatPacket{}
		if err := json.Unmarshal(payloadBytes, &chatPacket); err != nil {
			slog.Error("REPOSITORY", "unmarshal payload", err.Error())
			continue
		}

		chatHistory = append(chatHistory, &chatPacket)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return chatHistory, nil
}

func (s *scyllaRepository) GetLatestMessageFromBucket(ctx context.Context, conversation_id string) (*model.ChatPacket, error) {
	query := `
		SELECT payload
		FROM chat_history
		WHERE conversation_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	iter := s.session.Query(query, conversation_id).WithContext(ctx).Iter()

	var chatHistory []*model.ChatPacket
	var payloadBytes []byte

	for iter.Scan(&payloadBytes) {
		chatPacket := model.ChatPacket{}
		if err := json.Unmarshal(payloadBytes, &chatPacket); err != nil {
			slog.Error("REPOSITORY", "unmarshal payload", err.Error())
			continue
		}

		chatHistory = append(chatHistory, &chatPacket)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	if chatHistory[0] == nil {
		return nil, fmt.Errorf("no messages")
	}

	return chatHistory[0], nil
}
