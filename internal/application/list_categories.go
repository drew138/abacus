package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) ListCategories(ctx context.Context) ([]domain.Category, error) {
	categories, err := a.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing categories: %w", err)
	}
	return categories, nil
}
