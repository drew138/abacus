package application

import (
	"context"
	"fmt"
)

func (a *Application) ListCategoryTotals(ctx context.Context) (map[string]map[string]int64, error) {
	totals, err := a.repo.ListCategoryTotals(ctx)
	if err != nil {
		return nil, fmt.Errorf("totalling categories: %w", err)
	}
	return totals, nil
}
