package database

import (
	"context"
	"testing"

	"github.com/shanto-323/Chat-Server-1/client-service/internal/database/model"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func MockDb() (string, error) {
	ctx := context.WithoutCancel(context.Background())
	psContainer, err := postgres.Run(
		ctx,
		"postgres:latest",
		postgres.WithDatabase("mock"),
		postgres.WithUsername("shanto"),
		postgres.WithPassword("shanto323"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return "", err
	}

	connStr, err := psContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return "", nil
	}

	return connStr, nil
}

func TestCreateDatabase(t *testing.T) {
	t.Parallel()

	connStr, err := MockDb()
	assert.Nil(t, err)

	repo, err := NewUserRepository(connStr)

	assert.Nil(t, err)
	assert.NotNil(t, repo)
}

func TestInsertUser(t *testing.T) {
	t.Parallel()

	connStr, _ := MockDb()
	repo, _ := NewUserRepository(connStr)

	err := repo.InsertUser(context.Background(), &model.User{
		Username: "mock_user_u",
		Password: "mock_user_p",
	})
	assert.Nil(t, err)
}

func TestGetAccount(t *testing.T) {
	connStr, _ := MockDb()
	repo, _ := NewUserRepository(connStr)

	_ = repo.InsertUser(context.Background(), &model.User{
		Username: "mock_user_u",
		Password: "mock_user_p",
	})

	resp, err := repo.GetUser(context.Background(), "mock_user_u")
	assert.Nil(t, err)
	assert.Equal(t, "mock_user_p", resp.Password)
}

func TestDeleteAccount(t *testing.T) {
	connStr, _ := MockDb()
	repo, _ := NewUserRepository(connStr)

	_ = repo.InsertUser(context.Background(), &model.User{
		Username: "mock_user_u",
		Password: "mock_user_p",
	})

	err := repo.DeleteUser(context.Background(), "mock_user_u")
	assert.Nil(t, err)

	resp, _ := repo.GetUser(context.Background(), "mock_user_u")
	assert.Equal(t, uint(0), resp.ID)
}
