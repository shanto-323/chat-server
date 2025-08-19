package remote

import (
	"net/http"
	"time"
)

const (
	Get = "GET"
	Set = "SET"
)

type CacheClient interface{}

type cacheClient struct {
	baseUrl string
	client  *http.Client
}

func NewCacheClient() CacheClient {
	return &cacheClient{
		baseUrl: "http://client-service:8081/api/v1/client.service",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
