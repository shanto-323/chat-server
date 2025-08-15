package model

import "encoding/json"

type MessagePacket struct {
	MsgType    string          `json:"type"`
	SenderId   uint            `json:"sender_id,omitempty"`
	ReceiverId uint            `json:"receiver_id,omitempty"`
	Payload    json.RawMessage `json:"payload,omitempty"`
	Timestamp  int64           `json:"timestamp,omitempty"`
}

type UserAction struct {
	Action   string `json:"action"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ActivePool struct {
	AliveList []string `json:"alive_list"`
}

type Client struct {
	ID uint `json:"id"`
}
