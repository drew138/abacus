package application

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) ImportMovements(ctx context.Context, movements []domain.Movement, categories []domain.Category) (domain.ImportResult, error) {
	if err := a.categories.Ensure(ctx, categories); err != nil {
		return domain.ImportResult{}, fmt.Errorf("ensuring %d categories: %w", len(categories), err)
	}

	rules, err := a.movements.MerchantRules(ctx)
	if err != nil {
		return domain.ImportResult{}, fmt.Errorf("loading merchant rules: %w", err)
	}

	for i, m := range movements {
		movements[i] = applyImportRule(m, rules)
		movements[i].Source = domain.SourceStatement
	}

	result, err := a.movements.Import(ctx, movements)
	if err != nil {
		return domain.ImportResult{}, fmt.Errorf("importing %d movements: %w", len(movements), err)
	}
	return result, nil
}

func applyImportRule(m domain.Movement, rules map[string]string) domain.Movement {
	if m.IsSavings() {
		return m
	}
	if category, ok := rules[m.MerchantKey()]; ok {
		m.Category = category
	}
	return m
}
