package model

type CacheResponseAll struct {
	Status  int `json:"status"`
	Message struct {
		ConnPool []string `jons:"conn_pool"`
	} `json:"message"`
}
