package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shanto-323/Chat-Server-1/gateway-1/data/remote/model"
)

type CacheClient interface {
	AddActiveUser(connRequest *model.ConnRequest) error
	RemoveActiveUser(connRequest *model.ConnRequest) error
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

func (c *cacheClient) AddActiveUser(connRequest *model.ConnRequest) error {
	url := fmt.Sprintf("%s/cache/client.up", c.baseUrl)
	body, err := json.Marshal(&connRequest)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("err:%s :%d", err, resp.StatusCode)
	}

	return nil
}

func (c *cacheClient) RemoveActiveUser(connRequest *model.ConnRequest) error {
	url := fmt.Sprintf("%s/cache/client.close", c.baseUrl)
	body, err := json.Marshal(&connRequest)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("err:%s :%d", err, resp.StatusCode)
	}

	return nil
}
