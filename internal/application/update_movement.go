package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) UpdateMovement(ctx context.Context, movement domain.Movement, applyToMerchant bool) (domain.Movement, error) {
	current, err := a.repo.GetMovement(ctx, movement.ID)
	if err != nil {
		return domain.Movement{}, fmt.Errorf("finding movement %d: %w", movement.ID, err)
	}

	categories, err := a.repo.ListCategories(ctx)
	if err != nil {
		return domain.Movement{}, fmt.Errorf("listing categories: %w", err)
	}

	if err := validateUpdateMovement(movement, categories); err != nil {
		return domain.Movement{}, err
	}

	movement.Ref = current.Ref
	movement.Source = current.Source
	movement.Deleted = current.Deleted
	movement.Edited = current.Source == domain.SourceStatement

	saved, err := a.repo.UpdateMovement(ctx, movement)
	if err != nil {
		return domain.Movement{}, fmt.Errorf("saving movement %d: %w", movement.ID, err)
	}

	if applyToMerchant {
		if _, err := a.repo.RecategoriseMerchant(ctx, saved.Merchant, saved.Category); err != nil {
			return domain.Movement{}, fmt.Errorf("recategorising the merchant of movement %d: %w", saved.ID, err)
		}
	}
	return saved, nil
}

func validateUpdateMovement(movement domain.Movement, categories []domain.Category) error {
	for _, c := range categories {
		if c.Name == movement.Category {
			return nil
		}
	}
	return domain.ErrUnknownCategory
}
