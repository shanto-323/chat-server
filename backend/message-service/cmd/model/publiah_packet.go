package model

type PublishPacket struct {
	SessionId string `json:"session_id"`
	Data      string `json:"data"`
}
