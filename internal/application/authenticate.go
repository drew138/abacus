package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/token"
)

func (a *Application) Authenticate(ctx context.Context, raw string) (domain.User, error) {
	if raw == "" {
		return domain.User{}, domain.ErrNotAuthenticated
	}

	userID, expiresAt, err := a.repo.FindSession(ctx, token.Hash(raw))
	if errors.Is(err, domain.ErrNotFound) {
		return domain.User{}, domain.ErrNotAuthenticated
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("finding session: %w", err)
	}
	if !a.now().Before(expiresAt) {
		return domain.User{}, domain.ErrNotAuthenticated
	}

	user, err := a.repo.FindUser(ctx, userID)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.User{}, domain.ErrNotAuthenticated
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("finding user %d: %w", userID, err)
	}
	return user, nil
}
