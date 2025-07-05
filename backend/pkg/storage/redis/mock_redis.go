package redis

import (
	"context"
)

type MockRedisClient struct{}

func NewMockRedisClient() RedisClient {
	return &MockRedisClient{}
}

func (m *MockRedisClient) Close() {}

func (m *MockRedisClient) Set(ctx context.Context, key string, value any) error {
	return nil
}

func (m *MockRedisClient) Exists(ctx context.Context, key string) error {
	return nil
}

func (m *MockRedisClient) Get(ctx context.Context, key string) error {
	return nil
}

func (m *MockRedisClient) Remove(ctx context.Context, keys ...string) error {
	return nil
}

func (m *MockRedisClient) SetList(ctx context.Context, key string, value any) error {
	return nil
}

func (m *MockRedisClient) GetList(ctx context.Context, key string) ([]string, error) {
	return []string{}, nil
}

func (m *MockRedisClient) RemoveFromList(ctx context.Context, key string, value any) error {
	return nil
}
