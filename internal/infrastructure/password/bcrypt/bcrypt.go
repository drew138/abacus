package bcrypt

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/drew138/abacus/internal/domain"
)

type Passwords struct{}

func New() Passwords {
	return Passwords{}
}

func (Passwords) Hash(_ context.Context, plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w: %w", domain.ErrInternal, err)
	}
	return string(hash), nil
}

func (Passwords) Matches(_ context.Context, hash, plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("comparing password: %w: %w", domain.ErrInternal, err)
	}
	return true, nil
}
