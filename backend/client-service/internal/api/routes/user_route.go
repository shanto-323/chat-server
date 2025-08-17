package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shanto-323/Chat-Server-1/client-service/internal/api/model"
	"github.com/shanto-323/Chat-Server-1/client-service/internal/database"
	"github.com/shanto-323/Chat-Server-1/client-service/util"
)

type UserRoute interface {
	SignUpHandler(w http.ResponseWriter, r *http.Request) error
	SignInHandler(w http.ResponseWriter, r *http.Request) error
	DeleteUserHandler(w http.ResponseWriter, r *http.Request) error
}

type userRouteHnadler struct {
	service *database.UserService
}

func NewUserRoute(s *database.UserService) UserRoute {
	return &userRouteHnadler{service: s}
}

func (u *userRouteHnadler) SignUpHandler(w http.ResponseWriter, r *http.Request) error {
	authRequest := model.AuthRequest{}
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := u.service.SignUp(ctx, authRequest.Username, authRequest.Password)
	if err != nil {
		return err
	}

	return util.WriteJson(w, 200, resp)
}

func (u *userRouteHnadler) SignInHandler(w http.ResponseWriter, r *http.Request) error {
	authRequest := model.AuthRequest{}
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := u.service.SignIn(ctx, authRequest.Username, authRequest.Password)
	if err != nil {
		return err
	}

	return util.WriteJson(w, 200, resp)
}

func (u *userRouteHnadler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) error {
	v := mux.Vars(r)
	username := v["username"]

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := u.service.DeleteUser(ctx, username)
	if err != nil {
		return err
	}

	return util.WriteJson(w, 200, nil)
}
