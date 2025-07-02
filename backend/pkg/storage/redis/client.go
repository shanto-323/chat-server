package redis

import (
	"context"
	"time"

	"chat_app/backend/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisClient struct {
	client *redis.Client
	logger *logger.ZapLogger
}

func NewRedisClient(url string, logger *logger.ZapLogger) (*RedisClient, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return &RedisClient{
		client: client,
		logger: logger,
	}, nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value any) error {
	err := r.client.Set(ctx, key, value, 24*3600*time.Second).Err() // 1Day default
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	return nil
}

func (r *RedisClient) Exists(ctx context.Context, key string) error {
	err := r.client.Exists(ctx, key).Err()
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	return nil
}

func (r *RedisClient) Get(ctx context.Context, key string) error {
	resp := r.client.Get(ctx, key)
	val, err := resp.Result()
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	r.logger.Info(val)
	return nil
}

func (r *RedisClient) Remove(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	return nil
}

func (r *RedisClient) SetList(ctx context.Context, key string, value any) error {
	err := r.client.RPush(ctx, key, value).Err()
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	return nil
}

func (r *RedisClient) GetList(ctx context.Context, key string) ([]string, error) {
	resp := r.client.LRange(ctx, key, 0, -1)
	if resp.Err() != nil {
		r.logger.Error("redis", zap.String("err", resp.Err().Error()))
		return nil, resp.Err()
	}
	return resp.Val(), nil
}

func (r *RedisClient) RemoveFromList(ctx context.Context, key string, value any) error {
	err := r.client.LRem(ctx, key, 1, value).Err()
	if err != nil {
		r.logger.Error("redis", zap.String("err", err.Error()))
		return err
	}
	return nil
}
