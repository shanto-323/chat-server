package model

import "encoding/json"

type MessagePacket struct {
	MsgType    string          `json:"type"`
	SenderId   string          `json:"sender_id,omitempty"`
	ReceiverId string          `json:"receiver_id,omitempty"`
	Payload    json.RawMessage `json:"payload,omitempty"`
	Timestamp  int64           `json:"timestamp,omitempty"`
}

type ActivePool struct {
	AliveList []string `json:"alive_list"`
}

type Client struct {
	ID string `json:"id"`
}
