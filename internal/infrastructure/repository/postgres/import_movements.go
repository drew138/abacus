package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) ImportMovements(ctx context.Context, movements []domain.Movement) (added, updated, kept int, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, 0, domainError(err, "starting import")
	}
	defer func() { _ = tx.Rollback() }()

	rules, err := importMovementsRules(ctx, tx)
	if err != nil {
		return 0, 0, 0, err
	}

	for _, m := range movements {
		row := mapImportMovementsRow(m, rules)
		var inserted bool
		err := tx.QueryRowContext(ctx,
			`INSERT INTO movements (ref, occurred_on, merchant, merchant_key, amount, category, method, source)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 ON CONFLICT (ref) DO UPDATE SET occurred_on = excluded.occurred_on, merchant = excluded.merchant, merchant_key = excluded.merchant_key,
			     amount = excluded.amount, category = excluded.category, method = excluded.method, updated_at = now()
			 WHERE movements.edited = false AND movements.deleted_at IS NULL
			   AND (movements.occurred_on, movements.merchant, movements.amount, movements.category, movements.method)
			       IS DISTINCT FROM (excluded.occurred_on, excluded.merchant, excluded.amount, excluded.category, excluded.method)
			 RETURNING xmax = 0`,
			row.Ref, row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, row.Source).Scan(&inserted)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			kept++
		case err != nil:
			return 0, 0, 0, domainError(err, "importing a movement")
		case inserted:
			added++
		default:
			updated++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, 0, domainError(err, "committing import")
	}
	return added, updated, kept, nil
}

func importMovementsRules(ctx context.Context, tx *sql.Tx) (map[string]string, error) {
	rows, err := tx.QueryContext(ctx, "SELECT merchant_key, category FROM merchant_rules")
	if err != nil {
		return nil, domainError(err, "selecting merchant rules")
	}
	defer rows.Close()

	rules := map[string]string{}
	for rows.Next() {
		var key, category string
		if err := rows.Scan(&key, &category); err != nil {
			return nil, domainError(err, "reading a merchant rule")
		}
		rules[key] = category
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading merchant rules")
	}
	return rules, nil
}
