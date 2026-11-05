package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) AuthenticateUser(ctx context.Context, token string) (domain.User, error) {
	if token == "" {
		return domain.User{}, domain.ErrNotAuthenticated
	}

	userID, err := a.repo.GetSession(ctx, token)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.User{}, domain.ErrNotAuthenticated
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("finding session: %w", err)
	}

	user, err := a.repo.GetUser(ctx, userID)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.User{}, domain.ErrNotAuthenticated
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("finding user %d: %w", userID, err)
	}
	return user, nil
}
