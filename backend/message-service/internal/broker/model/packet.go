package model

type Packet struct {
	SenderId   string `json:"senderId" validate:"required"`
	ReceiverId string `json:"receiverId" validate:"required"`
	Data       string `json:"data,omitempty"`
}
