package model

type ConsumePacket struct {
	SessionId string `json:"session_id"`
	Data      string `json:"data"`
}
