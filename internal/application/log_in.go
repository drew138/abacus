package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/token"
)

func (a *Application) LogIn(ctx context.Context, username, plain string) (string, time.Time, error) {
	user, err := a.repo.FindByUsername(ctx, username)
	if errors.Is(err, domain.ErrNotFound) {
		return "", time.Time{}, domain.ErrInvalidCredentials
	}
	if err != nil {
		return "", time.Time{}, fmt.Errorf("finding user: %w", err)
	}

	ok, err := a.passwordMatches(user.PasswordHash, plain)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("checking password of user %d: %w: %w", user.ID, domain.ErrInternal, err)
	}
	if !ok {
		return "", time.Time{}, domain.ErrInvalidCredentials
	}

	raw, err := a.newToken()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("creating session token: %w: %w", domain.ErrInternal, err)
	}
	expiresAt := a.now().Add(a.sessionTTL)
	if err := a.repo.CreateSession(ctx, token.Hash(raw), user.ID, expiresAt); err != nil {
		return "", time.Time{}, fmt.Errorf("saving session of user %d: %w", user.ID, err)
	}
	return raw, expiresAt, nil
}
