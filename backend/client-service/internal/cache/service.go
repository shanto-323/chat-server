package cache

import (
	"context"
	"fmt"
)

type RedisService struct {
	client RedisClient
}

func NewRedisService(client RedisClient) *RedisService {
	return &RedisService{
		client: client,
	}
}

func (r *RedisService) Close() {
	r.client.Close()
}

// USER_ID = uid
// HASH_KEY = hkey
// SESSION_ID = session_id
// GATEWAY_ID = gatekey
func (r *RedisService) AddActiveUser(ctx context.Context, uid, hkey, gateKey, session_id string) error {
	// CHECK IF USER ALREADY EXISTS
	exist, err := r.client.CheckUID(ctx, uid)
	if err != nil {
		return err
	}

	// USER NOT INSIDE ACTIVE POOL
	if !exist {
		if err := r.client.Insert(ctx, uid, hkey); err != nil {
			return err
		}
	}

	// THIS MEANS USER WITH ID(uid) HAS A DEVICE(session_id) CONNECTED IN THIS(gateway_id) SERVER.
	if err := r.client.HashSetInsert(ctx, uid, gateKey, session_id); err != nil {
		return err
	}

	return nil
}

func (r *RedisService) RemoveActiveUser(ctx context.Context, uid, session_id, hkey string) error {
	// CHECK IF USER ALREADY EXISTS
	exist, err := r.client.CheckUID(ctx, uid)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("uid not exists %s", uid)
	}

	// HashSetRemove FUNC WILL PRE CALCULATE THE LENGTH AND RETURN THE SIZE AND REMOVE THE VALUE
	// HASH_SIZE IS SIZE-1
	size, err := r.client.HashSetRemove(ctx, uid, hkey)
	if err != nil {
		return err
	}

	// MEANS THIS WAS THE LAST USER.
	if size == 1 {
		if err := r.client.Remove(ctx, uid); err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisService) GetActivePool(ctx context.Context, uid string) (map[string]string, error) {
	// CHECK IF USER ALREADY EXISTS
	exist, err := r.client.CheckUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("uid not exists %s", uid)
	}

	return r.client.HashSetGet(ctx, uid)
}
