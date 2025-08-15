package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shanto-323/Chat-Server-1/client-service/internal/cache"
	"github.com/shanto-323/Chat-Server-1/client-service/util"
)

type CacheRoute interface {
	AddConnectionHandler(w http.ResponseWriter, r *http.Request) error
	RemoveConnectionHandler(w http.ResponseWriter, r *http.Request) error
	CheckConnectionHandler(w http.ResponseWriter, r *http.Request) error
}

type cacheRouteHandler struct {
	cache cache.RedisClient
}

func NewCacheRoute(c cache.RedisClient) CacheRoute {
	return &cacheRouteHandler{cache: c}
}

func (c *cacheRouteHandler) AddConnectionHandler(w http.ResponseWriter, r *http.Request) error {
	v := mux.Vars(r)
	id := v["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := c.cache.AddConnection(ctx, id, "active"); err != nil {
		return err
	}

	return util.WriteJson(w, 200, nil)
}

func (c *cacheRouteHandler) RemoveConnectionHandler(w http.ResponseWriter, r *http.Request) error {
	v := mux.Vars(r)
	id := v["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := c.cache.RemoveConnection(ctx, id); err != nil {
		return err
	}

	return util.WriteJson(w, 200, nil)
}

func (c *cacheRouteHandler) CheckConnectionHandler(w http.ResponseWriter, r *http.Request) error {
	v := mux.Vars(r)
	id := v["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := c.cache.CheckConnection(ctx, id); err != nil {
		return err
	}

	return util.WriteJson(w, 200, nil)
}
