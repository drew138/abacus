package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) ImportMovements(ctx context.Context, movements []domain.Movement, categories []domain.Category) (added, updated, kept int, err error) {
	if err := a.repo.EnsureCategories(ctx, categories); err != nil {
		return 0, 0, 0, fmt.Errorf("ensuring %d categories: %w", len(categories), err)
	}

	for i := range movements {
		movements[i].Source = domain.SourceStatement
	}

	added, updated, kept, err = a.repo.ImportMovements(ctx, movements)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("importing %d movements: %w", len(movements), err)
	}
	return added, updated, kept, nil
}
