package application

import (
	"context"
	"fmt"
)

func (a *Application) LogOutUser(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	if err := a.repo.DeleteSession(ctx, token); err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}
	return nil
}
