package model

type ConnRequest struct {
	ID        string `json:"id"`
	SessionId string `json:"session_id"`
	GatewayId string `json:"gateway_id"`
}
