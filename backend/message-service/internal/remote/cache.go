package remote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shanto-323/Chat-Server-1/message-service/internal/remote/model"
)

type CacheClient interface {
	GetActivePool(id string) (*model.CacheResponse, error)
	GetAlConnPool() (*model.CacheResponseAll, error)
}

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

func (c *cacheClient) GetActivePool(id string) (*model.CacheResponse, error) {
	url := fmt.Sprintf("%s/cache/client.get/%s", c.baseUrl, id)

	resp, err := c.client.Get(url)
	if err != nil || resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("err:%s", string(body))
	}
	defer resp.Body.Close()

	cacheResponse := model.CacheResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&cacheResponse); err != nil {
		return nil, err
	}

	return &cacheResponse, nil
}

func (c *cacheClient) GetAlConnPool() (*model.CacheResponseAll, error) {
	url := fmt.Sprintf("%s/cache/client.get", c.baseUrl)

	resp, err := c.client.Get(url)
	if err != nil || resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("err:%s", string(body))
	}
	defer resp.Body.Close()

	CacheResponseAll := model.CacheResponseAll{}
	if err := json.NewDecoder(resp.Body).Decode(&CacheResponseAll); err != nil {
		return nil, err
	}

	return &CacheResponseAll, nil
}
