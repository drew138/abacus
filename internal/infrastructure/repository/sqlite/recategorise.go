package sqlite

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Recategorise(ctx context.Context, merchantKey, category string) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, domainError(err, "starting recategorisation")
	}
	defer func() { _ = tx.Rollback() }()

	at := now()
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO merchant_rules (merchant_key, category, updated_at) VALUES (?, ?, ?)
		 ON CONFLICT (merchant_key) DO UPDATE SET category = excluded.category, updated_at = excluded.updated_at`,
		merchantKey, category, at); err != nil {
		return 0, domainError(err, "saving a merchant rule")
	}
	res, err := tx.ExecContext(ctx,
		"UPDATE movements SET category = ?, updated_at = ? WHERE merchant_key = ? AND category <> ? AND category <> ?",
		category, at, merchantKey, category, domain.CategorySavings)
	if err != nil {
		return 0, domainError(err, "recategorising movements")
	}
	if err := tx.Commit(); err != nil {
		return 0, domainError(err, "committing recategorisation")
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}
