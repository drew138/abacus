package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

const minRegisterUserPasswordLength = 12

func (a *Application) RegisterUser(ctx context.Context, username, plain string) (domain.User, error) {
	users, err := a.repo.CountUsers(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("checking for users: %w", err)
	}

	if err := validateRegisterUser(users > 0, plain); err != nil {
		return domain.User{}, err
	}

	user, err := a.repo.CreateFirstUser(ctx, username, plain)
	if err != nil {
		return domain.User{}, fmt.Errorf("creating the first user: %w", err)
	}
	user.PasswordHash = ""
	return user, nil
}

func validateRegisterUser(taken bool, plain string) error {
	if taken {
		return domain.ErrRegistrationClosed
	}
	if len([]rune(plain)) < minRegisterUserPasswordLength {
		return domain.ErrPasswordTooShort
	}
	return nil
}
