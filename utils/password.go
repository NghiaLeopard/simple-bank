package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordCode,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)

	if err != nil {
		return "",fmt.Errorf("fail to hash password: %w",err)
	}

	return string(passwordCode),nil
}

func CheckPassword(password string, hasPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hasPassword),[]byte(password))
}