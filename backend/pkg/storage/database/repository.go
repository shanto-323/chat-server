package database

import (
	"chat_app/backend/pkg/model"

	"github.com/gocql/gocql"
)

type Repository interface {
	Close()
	UpsertChat(c model.Chat) error
	UpsertOffline(c model.OfflineChat) error
}

type scyllaDb struct {
	session *gocql.Session
}

func NewScyllaDb(connIp string) (Repository, error) {
	cluster := gocql.NewCluster(connIp)
	cluster.Port = 9042
	cluster.Keyspace = "cluster"
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &scyllaDb{
		session: session,
	}, nil
}

func (s *scyllaDb) Close() {
	s.session.Close()
}

func (s *scyllaDb) UpsertChat(c model.Chat) error {
	query := `
		INSERT INTO chat_history(
			chat_id,message,created_at
		) VALUES (?,?,?)
	`
	return s.session.Query(query, c.ChatId, c.Message, c.CreatedAt).Exec()
}

func (s *scyllaDb) UpsertOffline(c model.OfflineChat) error {
	query := `
		INSERT INTO chat_offline(
			sender_id,receiver_id,message,created_at
		) VALUES (?,?,?,?)
	`
	return s.session.Query(query, c.SenderId, c.ReceiverId, c.Message, c.CreatedAt).Exec()
}
