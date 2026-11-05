package application

import (
	"context"
	"fmt"
)

func (a *Application) CategoryTotals(ctx context.Context) (map[string]map[string]int64, error) {
	totals, err := a.repo.Totals(ctx)
	if err != nil {
		return nil, fmt.Errorf("totalling categories: %w", err)
	}
	return totals, nil
}
