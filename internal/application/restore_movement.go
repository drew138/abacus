package application

import (
	"context"
	"fmt"
)

func (a *Application) RestoreMovement(ctx context.Context, id int64) error {
	if err := a.repo.RestoreMovement(ctx, id); err != nil {
		return fmt.Errorf("restoring movement %d: %w", id, err)
	}
	return nil
}
