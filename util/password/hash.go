package password

import (
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(password string) (string, error) {
	var (
		err   error
		bytes []byte
	)
	copy(bytes[:], password)
	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}
