package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Import(ctx context.Context, movements []domain.Movement) (added, updated, kept int, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, 0, domainError(err, "starting import")
	}
	defer func() { _ = tx.Rollback() }()

	rules, err := importRules(ctx, tx)
	if err != nil {
		return 0, 0, 0, err
	}

	at := now()
	for _, m := range movements {
		row := mapMovementRow(m)
		if rule, ok := rules[row.MerchantKey]; ok && row.Category != domain.CategorySavings {
			row.Category = rule
		}
		var id int64
		var edited bool
		var deletedAt sql.NullString
		err := tx.QueryRowContext(ctx, "SELECT id, edited, deleted_at FROM movements WHERE ref = ?", row.Ref).Scan(&id, &edited, &deletedAt)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO movements (ref, occurred_on, merchant, merchant_key, amount, category, method, note, source, edited, created_at, updated_at)
				 VALUES (?, ?, ?, ?, ?, ?, ?, '', ?, 0, ?, ?)`,
				row.Ref, row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, row.Source, at, at); err != nil {
				return 0, 0, 0, domainError(err, "inserting an imported movement")
			}
			added++
		case err != nil:
			return 0, 0, 0, domainError(err, "finding an imported movement")
		case edited || deletedAt.Valid:
			kept++
		default:
			res, err := tx.ExecContext(ctx,
				`UPDATE movements SET occurred_on = ?, merchant = ?, merchant_key = ?, amount = ?, category = ?, method = ?, updated_at = ?
				 WHERE id = ? AND (occurred_on <> ? OR merchant <> ? OR amount <> ? OR category <> ? OR method <> ?)`,
				row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, at,
				id, row.OccurredOn, row.Merchant, row.Amount, row.Category, row.Method)
			if err != nil {
				return 0, 0, 0, domainError(err, fmt.Sprintf("updating imported movement %d", id))
			}
			if n, _ := res.RowsAffected(); n > 0 {
				updated++
			} else {
				kept++
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, 0, domainError(err, "committing import")
	}
	return added, updated, kept, nil
}

func importRules(ctx context.Context, tx *sql.Tx) (map[string]string, error) {
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
