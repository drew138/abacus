package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Import(ctx context.Context, movements []domain.Movement) (domain.ImportResult, error) {
	var result domain.ImportResult
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return result, domainError(err, "starting import")
	}
	defer func() { _ = tx.Rollback() }()

	at := now()
	for _, m := range movements {
		row := mapMovementRow(m)
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
				return domain.ImportResult{}, domainError(err, "inserting an imported movement")
			}
			result.Added++
		case err != nil:
			return domain.ImportResult{}, domainError(err, "finding an imported movement")
		case edited || deletedAt.Valid:
			result.Kept++
		default:
			res, err := tx.ExecContext(ctx,
				`UPDATE movements SET occurred_on = ?, merchant = ?, merchant_key = ?, amount = ?, category = ?, method = ?, updated_at = ?
				 WHERE id = ? AND (occurred_on <> ? OR merchant <> ? OR amount <> ? OR category <> ? OR method <> ?)`,
				row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, at,
				id, row.OccurredOn, row.Merchant, row.Amount, row.Category, row.Method)
			if err != nil {
				return domain.ImportResult{}, domainError(err, fmt.Sprintf("updating imported movement %d", id))
			}
			if n, _ := res.RowsAffected(); n > 0 {
				result.Updated++
			} else {
				result.Kept++
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return domain.ImportResult{}, domainError(err, "committing import")
	}
	return result, nil
}
