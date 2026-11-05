package postgres

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) EnsureCategories(ctx context.Context, categories []domain.Category) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domainError(err, "starting to ensure categories")
	}
	defer func() { _ = tx.Rollback() }()

	for _, c := range categories {
		row := mapEnsureCategoriesRow(c)
		if _, err := tx.ExecContext(ctx, "INSERT INTO categories (name, parent) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING", row.Name, row.Parent); err != nil {
			return domainError(err, "inserting a category")
		}
	}
	if err := tx.Commit(); err != nil {
		return domainError(err, "committing categories")
	}
	return nil
}
