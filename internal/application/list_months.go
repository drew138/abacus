package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) ListMonths(ctx context.Context) ([]domain.MonthSummary, error) {
	months, err := a.repo.Months(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing months: %w", err)
	}
	return months, nil
}
