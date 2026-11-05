package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

const minPasswordLength = 12

func (a *Application) SetPassword(ctx context.Context, username, plain string) (domain.User, error) {
	if err := validateSetPassword(plain); err != nil {
		return domain.User{}, err
	}

	hash, err := a.passwords.Hash(ctx, plain)
	if err != nil {
		return domain.User{}, fmt.Errorf("hashing password: %w", err)
	}

	user, err := a.users.SavePassword(ctx, username, hash)
	if err != nil {
		return domain.User{}, fmt.Errorf("saving password: %w", err)
	}

	if err := a.sessions.DeleteSessionsForUser(ctx, user.ID); err != nil {
		return domain.User{}, fmt.Errorf("signing out user %d: %w", user.ID, err)
	}
	user.PasswordHash = ""
	return user, nil
}

func validateSetPassword(plain string) error {
	if len([]rune(plain)) < minPasswordLength {
		return domain.ErrPasswordTooShort
	}
	return nil
}
