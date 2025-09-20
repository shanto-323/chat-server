package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shanto-323/Chat-Server-1/gateway-1/data/remote/model"
)

const (
	SignUP = "SIGNUP"
	SignIN = "SIGNIN"
)

type UserClient interface {
	Auth(payload []byte) (*model.User, error)
}

type client struct {
	baseUrl string
	client  *http.Client
}

func NewClient() UserClient {
	return &client{
		baseUrl: "http://client-service:8081/api/v1/client.service",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (u *client) Auth(payload []byte) (*model.User, error) {
	request := model.UserRequest{}
	if err := json.Unmarshal(payload, &request); err != nil {
		return nil, err
	}

	method := strings.ToUpper(strings.ReplaceAll(request.Method, " ", ""))
	switch method {
	case SignUP:
		return u.singUp(request.Username, request.Password)
	case SignIN:
		return u.singIn(request.Username, request.Password)
	}

	return nil, fmt.Errorf("unknown method")
}

func (u *client) singUp(username, password string) (*model.User, error) {
	body, err := json.Marshal(model.UserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/user/sign.up", u.baseUrl)

	resp, err := u.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		var signInResponse model.SignUpResponse
		if err := json.NewDecoder(resp.Body).Decode(&signInResponse); err != nil {
			return nil, fmt.Errorf("error unmarshal data")
		}
		return &signInResponse.Message, nil
	}

	return nil, fmt.Errorf("denied")
}

func (u *client) singIn(username, password string) (*model.User, error) {
	body, err := json.Marshal(model.UserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/user/sign.in", u.baseUrl)

	resp, err := u.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		var signInResponse model.SignUpResponse
		if err := json.NewDecoder(resp.Body).Decode(&signInResponse); err != nil {
			return nil, fmt.Errorf("error unmarshal data")
		}
		return &signInResponse.Message, nil
	}

	return nil, fmt.Errorf("denied")
}
