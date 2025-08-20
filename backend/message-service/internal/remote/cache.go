package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shanto-323/Chat-Server-1/message-service/internal/remote/model"
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

func (c *cacheClient) GetActivePool(connRequest *model.ConnRequest) (*model.CacheResponse, error) {
	url := fmt.Sprintf("%s/cache/client.get", c.baseUrl)
	body, err := json.Marshal(&connRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("err:%s :%d", err, resp.StatusCode)
	}

	cacheResponse := model.CacheResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&cacheResponse); err != nil {
		return nil, err
	}

	return &cacheResponse, nil
}
