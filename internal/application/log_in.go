package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/token"
)

func (a *Application) LogIn(ctx context.Context, username, plain string) (string, domain.Session, error) {
	user, err := a.repo.FindByUsername(ctx, username)
	if errors.Is(err, domain.ErrNotFound) {
		return "", domain.Session{}, domain.ErrInvalidCredentials
	}
	if err != nil {
		return "", domain.Session{}, fmt.Errorf("finding user: %w", err)
	}

	ok, err := a.passwordMatches(user.PasswordHash, plain)
	if err != nil {
		return "", domain.Session{}, fmt.Errorf("checking password of user %d: %w: %w", user.ID, domain.ErrInternal, err)
	}
	if !ok {
		return "", domain.Session{}, domain.ErrInvalidCredentials
	}

	raw, err := a.newToken()
	if err != nil {
		return "", domain.Session{}, fmt.Errorf("creating session token: %w: %w", domain.ErrInternal, err)
	}
	session := domain.Session{TokenHash: token.Hash(raw), UserID: user.ID, ExpiresAt: a.now().Add(a.sessionTTL)}

	if err := a.repo.CreateSession(ctx, session); err != nil {
		return "", domain.Session{}, fmt.Errorf("saving session of user %d: %w", user.ID, err)
	}
	return raw, session, nil
}
