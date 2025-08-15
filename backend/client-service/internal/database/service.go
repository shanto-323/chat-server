package database

import (
	"context"

	"github.com/shanto-323/Chat-Server-1/client-service/internal/database/model"
	"github.com/shanto-323/Chat-Server-1/client-service/util/hash"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) SignUp(ctx context.Context, username, password string) (*model.User, error) {
	hash, err := hash.GenerateHash(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: hash,
	}

	if err := u.repo.InsertUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) SignIn(ctx context.Context, username, password string) (*model.User, error) {
	resp, err := u.repo.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	if err := hash.CompareWithHash([]byte(resp.Password), []byte(password)); err != nil {
		return nil, err
	}

	return resp, nil
}

func (u *UserService) DeleteUser(ctx context.Context, username string) error {
	return u.repo.DeleteUser(ctx, username)
}
