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

	session, err := a.repo.FindSession(ctx, token.Hash(raw))
	if errors.Is(err, domain.ErrNotFound) {
		return domain.User{}, domain.ErrNotAuthenticated
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("finding session: %w", err)
	}
	if session.Expired(a.now()) {
		return domain.User{}, domain.ErrNotAuthenticated
	}

	user, err := a.repo.FindUser(ctx, session.UserID)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.User{}, domain.ErrNotAuthenticated
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("finding user %d: %w", session.UserID, err)
	}
	return user, nil
}
