package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const (
	LIST_KEY = "list" // Set Name.
)

type RedisClient interface {
	Close()
	Insert(ctx context.Context, uid, hkey string) error
	Remove(ctx context.Context, uid string) error
	CheckUID(ctx context.Context, uid string) (bool, error)
	HashSetInsert(ctx context.Context, uid, gateway, session_id string) error
	HashSetRemove(ctx context.Context, uid string, hashKey string) (int64, error)
	HashSetGet(ctx context.Context, uid string) (map[string]string, error)
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

// USER_ID = uid
// HASH_KEY = hkey
// SESSION_ID = session_id
// GATEWAY_ID = gatekey
func (r *redisClient) Insert(ctx context.Context, uid, hkey string) error {
	_, err := r.client.Set(ctx, uid, hkey, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisClient) Remove(ctx context.Context, uid string) error {
	_, err := r.client.Del(ctx, uid).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisClient) CheckUID(ctx context.Context, uid string) (bool, error) {
	resp, err := r.client.Exists(ctx, uid).Result()
	if err != nil {
		return false, err
	}
	return resp > 0, nil
}

func (r *redisClient) HashSetInsert(ctx context.Context, uid, gateway, session_id string) error {
	_, err := r.client.HSet(ctx, uid, gateway, session_id).Result()
	return err
}

func (r *redisClient) HashSetRemove(ctx context.Context, session_id, hashKey string) (int64, error) {
	size, err := r.client.HLen(ctx, hashKey).Result()
	if err != nil {
		return 0, nil
	}
	_, err = r.client.HDel(ctx, session_id, hashKey).Result()
	return size, err
}

func (r *redisClient) HashSetGet(ctx context.Context, uid string) (map[string]string, error) {
	return r.client.HGetAll(ctx, uid).Result()
}
