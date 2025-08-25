package model

type UserRequest struct {
	Method   string `json:"method" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
