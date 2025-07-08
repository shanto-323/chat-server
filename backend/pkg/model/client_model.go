package model

import "time"

type ClientModel struct {
	Id        string    `json:"id"`
	Conn      string    `json:"conn"`
	Alive     bool      `json:"alive"`
	LastAlive time.Time `json:"last_alive"`
}
