package hash

import "golang.org/x/crypto/bcrypt"

func CompareWithHash(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
