package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

const minRegisterPasswordLength = 12

func (a *Application) Register(ctx context.Context, username, plain string) (domain.User, error) {
	taken, err := a.users.HasUsers(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("checking for users: %w", err)
	}

	if err := validateRegister(taken, plain); err != nil {
		return domain.User{}, err
	}

	hash, err := a.passwords.Hash(ctx, plain)
	if err != nil {
		return domain.User{}, fmt.Errorf("hashing password: %w", err)
	}

	user, err := a.users.CreateFirstUser(ctx, username, hash)
	if err != nil {
		return domain.User{}, fmt.Errorf("creating the first user: %w", err)
	}
	user.PasswordHash = ""
	return user, nil
}

func validateRegister(taken bool, plain string) error {
	if taken {
		return domain.ErrRegistrationClosed
	}
	if len([]rune(plain)) < minRegisterPasswordLength {
		return domain.ErrPasswordTooShort
	}
	return nil
}
