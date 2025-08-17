package model

type SignUpResponse struct {
	Status  int  `json:"status"`
	Message User `json:"message"`
}
