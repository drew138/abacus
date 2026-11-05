package sqlite

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Ensure(ctx context.Context, categories []domain.Category) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domainError(err, "starting to ensure categories")
	}
	defer func() { _ = tx.Rollback() }()

	at := now()
	for _, c := range categories {
		row := mapCategoryRow(c)
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO categories (name, parent, created_at) VALUES (?, ?, ?) ON CONFLICT (name) DO NOTHING",
			row.Name, row.Parent, at); err != nil {
			return domainError(err, "inserting a category")
		}
	}
	if err := tx.Commit(); err != nil {
		return domainError(err, "committing categories")
	}
	return nil
}
