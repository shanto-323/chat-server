package redis

import (
	"context"
	"time"

	"chat_app/backend/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	LIST_KEY = "list"
)

type RedisClient interface {
	Close()
	Set(ctx context.Context, key string, value any) error
	Exists(ctx context.Context, key string) error
	Get(ctx context.Context, key string) error
	Remove(ctx context.Context, keys ...string) error
	SetList(ctx context.Context, key string, value any) error
	GetList(ctx context.Context, key string) ([]string, error)
	RemoveFromList(ctx context.Context, key string, value any) error
}

type redisClient struct {
	client *redis.Client
	logger logger.ZapLogger
}

func NewRedisClient(url string, logger logger.ZapLogger) (RedisClient, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return &redisClient{
		client: client,
		logger: logger,
	}, nil
}

func (r *redisClient) Close() {
	r.client.Close()
}

func (r *redisClient) Set(ctx context.Context, key string, value any) error {
	err := r.client.Set(ctx, key, value, 24*3600*time.Second).Err() // 1Day default
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	return nil
}

func (r *redisClient) Exists(ctx context.Context, key string) error {
	err := r.client.Exists(ctx, key).Err()
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	return nil
}

func (r *redisClient) Get(ctx context.Context, key string) error {
	resp := r.client.Get(ctx, key)
	val, err := resp.Result()
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	r.logger.Info(val)
	return nil
}

func (r *redisClient) Remove(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	return nil
}

func (r *redisClient) SetList(ctx context.Context, key string, value any) error {
	err := r.client.RPush(ctx, key, value).Err()
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	return nil
}

func (r *redisClient) GetList(ctx context.Context, key string) ([]string, error) {
	resp := r.client.LRange(ctx, key, 0, -1)
	if resp.Err() != nil {
		r.logger.Error("redis", zap.String("err", resp.Err().Error()))
		return nil, resp.Err()
	}
	return resp.Val(), nil
}

func (r *redisClient) RemoveFromList(ctx context.Context, key string, value any) error {
	err := r.client.LRem(ctx, key, 1, value).Err()
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	return nil
}
