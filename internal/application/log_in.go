package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/password"
)

const sessionLifetime = 30 * 24 * time.Hour

func (a *Application) LogIn(ctx context.Context, username, plain string) (string, time.Time, error) {
	user, err := a.repo.FindByUsername(ctx, username)
	if errors.Is(err, domain.ErrNotFound) {
		return "", time.Time{}, domain.ErrInvalidCredentials
	}
	if err != nil {
		return "", time.Time{}, fmt.Errorf("finding user: %w", err)
	}

	ok, err := password.Matches(user.PasswordHash, plain)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("checking password of user %d: %w: %w", user.ID, domain.ErrInternal, err)
	}
	if !ok {
		return "", time.Time{}, domain.ErrInvalidCredentials
	}

	token, expiresAt, err := a.repo.CreateSession(ctx, user.ID, sessionLifetime)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("creating a session for user %d: %w", user.ID, err)
	}
	return token, expiresAt, nil
}
