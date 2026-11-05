package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) RecategoriseMerchant(ctx context.Context, merchant, category string) (int, error) {
	categories, err := a.repo.List(ctx)
	if err != nil {
		return 0, fmt.Errorf("listing categories: %w", err)
	}

	if err := validateRecategoriseMerchant(category, categories); err != nil {
		return 0, err
	}

	changed, err := a.repo.Recategorise(ctx, merchant, category)
	if err != nil {
		return 0, fmt.Errorf("recategorising a merchant: %w", err)
	}
	return changed, nil
}

func validateRecategoriseMerchant(category string, categories []domain.Category) error {
	if category == domain.CategorySavings {
		return domain.ErrUnknownCategory
	}
	for _, c := range categories {
		if c.Name == category {
			return nil
		}
	}
	return domain.ErrUnknownCategory
}
