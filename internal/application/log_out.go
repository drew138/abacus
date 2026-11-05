package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/pkg/token"
)

func (a *Application) LogOut(ctx context.Context, raw string) error {
	if raw == "" {
		return nil
	}
	if err := a.repo.DeleteSession(ctx, token.Hash(raw)); err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}
	return nil
}
