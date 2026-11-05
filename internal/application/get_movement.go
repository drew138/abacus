package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) GetMovement(ctx context.Context, id int64) (domain.Movement, error) {
	movement, err := a.movements.Find(ctx, id)
	if err != nil {
		return domain.Movement{}, fmt.Errorf("finding movement %d: %w", id, err)
	}
	return movement, nil
}
