package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const (
	LIST_KEY = "list" // Set Name.
)

type RedisClient interface {
	Close()
	SAdd(ctx context.Context, pool ...any) error
	SRem(ctx context.Context, members ...any) error
	SMembers(ctx context.Context) ([]string, error)
	IsMember(ctx context.Context, key string) (bool, error)
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(url string) (RedisClient, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return &redisClient{
		client: client,
	}, nil
}

func (r *redisClient) Close() {
	r.client.Close()
}

func (r *redisClient) SAdd(ctx context.Context, pool ...any) error {
	_, err := r.client.SAdd(ctx, LIST_KEY, pool...).Result()
	return err
}

func (r *redisClient) SRem(ctx context.Context, members ...any) error {
	_, err := r.client.SRem(ctx, LIST_KEY, members...).Result()
	return err
}

func (r *redisClient) SMembers(ctx context.Context) ([]string, error) {
	resp, err := r.client.SMembers(ctx, LIST_KEY).Result()
	return resp, err
}

func (r *redisClient) IsMember(ctx context.Context, key string) (bool, error) {
	resp, err := r.client.SIsMember(ctx, LIST_KEY, key).Result()
	return resp, err
}
