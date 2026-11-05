package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) ListMovements(ctx context.Context, month string) ([]domain.Movement, error) {
	movements, err := a.repo.ListMovements(ctx, month)
	if err != nil {
		return nil, fmt.Errorf("listing movements for %s: %w", month, err)
	}
	return movements, nil
}
