package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	connStr, _ := MockDb()
	repo, _ := NewUserRepository(connStr)
	service := NewUserService(repo)

	_, err := service.SignUp(context.Background(), "mock_user_u", "mock_user_p")
	assert.Nil(t, err)

	_, err = service.SignIn(context.Background(), "mock_user_u", "mock_user_p")
	assert.Nil(t, err)

	_, err = service.SignIn(context.Background(), "mock_user_u", "mock_user_err")
	assert.NotNil(t, err)

	err = service.DeleteUser(context.Background(), "mock_user_u")
	assert.Nil(t, err)

	_, err = service.SignIn(context.Background(), "mock_user_u", "mock_user_p")
	assert.NotNil(t, err)
}
