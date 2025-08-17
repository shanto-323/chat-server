package model

const (
	TYPE_CHAT  = "chat"
	TYPE_LIST  = "list"
	TYPE_ALIVE = "alive"
)

type Packet struct {
	Type       string `json:"type" validate:"required"`
	SenderId   string `json:"senderId" validate:"required"`
	ReceiverId string `json:"receiverId" validate:"required"`
	Data       string `json:"data,omitempty"`
}
