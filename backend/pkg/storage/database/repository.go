package database

import (
	"chat_app/backend/pkg/model"

	"github.com/gocql/gocql"
)

type Repository interface {
	Close()
}

type scyllaDb struct {
	session *gocql.Session
}

func NewScyllaDb(connIp string) (Repository, error) {
	cluster := gocql.NewCluster(connIp)
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

func (s *scyllaDb) UpsertClient(c model.ClientModel) error {
	query := `
		INSERT INTO clients(
			id,conn,alive,last_alive
		) VALUES (?,?,?,?)
	`
	return s.session.Query(query, c.Id, c.Conn, c.Alive, c.LastAlive).Exec()
}
