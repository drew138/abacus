package application

import (
	"context"
	"fmt"
)

func (a *Application) DeleteMovement(ctx context.Context, id int64) error {
	if err := a.movements.SetDeleted(ctx, id, true); err != nil {
		return fmt.Errorf("deleting movement %d: %w", id, err)
	}
	return nil
}
