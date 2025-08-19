package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/shanto-323/Chat-Server-1/client-service/internal/api/model"
	"github.com/shanto-323/Chat-Server-1/client-service/internal/cache"
	"github.com/shanto-323/Chat-Server-1/client-service/util"
)

type CacheRoute interface {
	AddConnectionHandler(w http.ResponseWriter, r *http.Request) error
	RemoveConnectionHandler(w http.ResponseWriter, r *http.Request) error
	GetConnectionHandler(w http.ResponseWriter, r *http.Request) error
}

type cacheRouteHandler struct {
	cache *cache.RedisService
}

func NewCacheRoute(c *cache.RedisService) CacheRoute {
	return &cacheRouteHandler{cache: c}
}

// USER_ID = uid
// HASH_KEY = hkey = uid(for now)
// SESSION_ID = session_id
// GATEWAY_ID = gatekey
func (c *cacheRouteHandler) AddConnectionHandler(w http.ResponseWriter, r *http.Request) error {
	connRequest := model.ConnRequest{}
	if err := json.NewDecoder(r.Body).Decode(&connRequest); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := c.cache.AddActiveUser(ctx, connRequest.ID, connRequest.ID, connRequest.GatewayId, connRequest.SessionId); err != nil {
		return err
	}

	return util.WriteJson(w, 200, nil)
}

func (c *cacheRouteHandler) RemoveConnectionHandler(w http.ResponseWriter, r *http.Request) error {
	connRequest := model.ConnRequest{}
	if err := json.NewDecoder(r.Body).Decode(&connRequest); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := c.cache.RemoveActiveUser(ctx, connRequest.ID, connRequest.SessionId, connRequest.ID); err != nil {
		return err
	}

	return util.WriteJson(w, 200, nil)
}

func (c *cacheRouteHandler) GetConnectionHandler(w http.ResponseWriter, r *http.Request) error {
	connRequest := model.ConnRequest{}
	if err := json.NewDecoder(r.Body).Decode(&connRequest); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	activePool, err := c.cache.GetActivePool(ctx, connRequest.ID)
	if err != nil {
		return err
	}

	connResponse := model.ConnResponse{
		ActivePool: activePool,
	}

	return util.WriteJson(w, 200, connResponse)
}
