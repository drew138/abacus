package postgres

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) RecategoriseMerchant(ctx context.Context, merchant, category string) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, domainError(err, "starting recategorisation")
	}
	defer func() { _ = tx.Rollback() }()

	key := merchantKey(merchant)
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO merchant_rules (merchant_key, category) VALUES ($1, $2)
		 ON CONFLICT (merchant_key) DO UPDATE SET category = excluded.category, updated_at = now()`,
		key, category); err != nil {
		return 0, domainError(err, "saving a merchant rule")
	}
	res, err := tx.ExecContext(ctx,
		"UPDATE movements SET category = $1, updated_at = now() WHERE merchant_key = $2 AND category <> $1 AND category <> $3",
		category, key, domain.CategorySavings)
	if err != nil {
		return 0, domainError(err, "recategorising movements")
	}
	if err := tx.Commit(); err != nil {
		return 0, domainError(err, "committing recategorisation")
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}
