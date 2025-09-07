package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func MockRedisDb() (string, error) {
	ctx := context.WithoutCancel(context.Background())
	psContainer, err := redis.Run(
		ctx,
		"redis:8.2-m01-alpine3.22",
	)
	if err != nil {
		return "", err
	}

	connStr, err := psContainer.ConnectionString(ctx)
	if err != nil {
		return "", nil
	}

	return connStr, nil
}

func TestRedisClient(t *testing.T) {
	connStr, err := MockRedisDb()
	assert.Nil(t, err)

	client, err := NewRedisClient(connStr)
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestAddActiveUser(t *testing.T) {
	connStr, _ := MockRedisDb()
	client, _ := NewRedisClient(connStr)
	service := NewRedisService(client)
	defer service.Close()

	err := service.AddActiveUser(context.Background(), "mock_uid", "mock_hkey", "mock_gateKey", "mock_session_id")
	assert.Nil(t, err)

	resp, err := service.GetActivePool(context.Background(), "mock_uid")
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	resps, err := service.GetActiveUsers(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, resps)

	err = service.RemoveActiveUser(context.Background(), "mock_uid", "mock_session_id", "mock_gateKey")
	assert.Nil(t, err)
}
