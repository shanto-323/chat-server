package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	password := "mock_user_p"

	hash, err := GenerateHash(password)
	assert.Nil(t, err)

	err = CompareWithHash(hash, password)
	assert.Nil(t, err)
}
