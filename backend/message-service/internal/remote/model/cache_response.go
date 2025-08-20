package model

type CacheResponse struct {
	Status  int `json:"status"`
	Message struct {
		ActivePool map[string]string `json:"active_pool"`
	} `json:"message"`
}
