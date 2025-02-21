package hashingapp

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type HashMngr struct {
	secretKey string
}

func NewHashMngr(secretKey string) *HashMngr {
	return &HashMngr{
		secretKey: secretKey,
	}
}

func (obj *HashMngr) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password+obj.secretKey), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hashed), nil
}

func (obj *HashMngr) CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+obj.secretKey))
}
