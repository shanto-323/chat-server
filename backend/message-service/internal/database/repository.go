package database

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/database/model"
)

type MessageRepository interface {
	Close()
	InsertMessage(ctx context.Context, chat *model.Chat) error
	GetMessageFromBucket(ctx context.Context, conversation_id string, createdAt time.Time) ([]string, error)
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
            chat_id, conversation_id, sender_id, receiver_id, message, created_at, offline
        ) VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	return s.session.Query(query,
		gocql.TimeUUID(),
		chat.ConversationID,
		chat.SenderID,
		chat.ReceiverID,
		chat.Message,
		chat.CreatedAt,
		chat.Offline,
	).WithContext(ctx).Exec()
}

func (s *scyllaRepository) GetMessageFromBucket(ctx context.Context, conversation_id string, createdAt time.Time) ([]string, error) {
	query := `
		SELECT message
		FROM chat_history
		WHERE conversation_id = ?
		  AND created < ?
		ORDER BY created DESC
		LIMIT 10
	`
	iter := s.session.Query(query, conversation_id, createdAt).WithContext(ctx).Iter()

	var messages []string
	var msg string
	for iter.Scan(&msg) {
		messages = append(messages, msg)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return messages, nil
}
