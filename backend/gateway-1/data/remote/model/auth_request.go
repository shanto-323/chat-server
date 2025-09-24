package model

import "encoding/json"

type UserRequest struct {
	Method     string          `json:"method" validate:"required"`
	Credential json.RawMessage `json:"credentials" validate:"required"`
}

type MethodSignUp struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type MethodSignIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
