package hash

import "golang.org/x/crypto/bcrypt"

func CompareWithHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
