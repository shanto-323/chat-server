package model

type Message struct {
	UserId   string `json:"userId" validate:"required"`
	SenderId string `json:"senderId" validate:"required"`
	Messsage string `json:"message" validate:"required"`
}
