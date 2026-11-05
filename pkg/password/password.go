package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hash(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(hash), nil
}

func Matches(hash, plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("comparing password: %w", err)
	}
	return true, nil
}
