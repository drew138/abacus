package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) CreateMovement(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	categories, err := a.repo.ListCategories(ctx)
	if err != nil {
		return domain.Movement{}, fmt.Errorf("listing categories: %w", err)
	}

	if err := validateCreateMovement(movement, categories); err != nil {
		return domain.Movement{}, err
	}

	movement.Ref = ""
	movement.Source = domain.SourceManual
	movement.Edited = false
	movement.Deleted = false

	saved, err := a.repo.CreateMovement(ctx, movement)
	if err != nil {
		return domain.Movement{}, fmt.Errorf("saving movement: %w", err)
	}
	return saved, nil
}

func validateCreateMovement(movement domain.Movement, categories []domain.Category) error {
	for _, c := range categories {
		if c.Name == movement.Category {
			return nil
		}
	}
	return domain.ErrUnknownCategory
}
