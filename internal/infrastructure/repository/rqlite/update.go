package rqlite

import (
	"context"
	"fmt"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Update(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	row := mapUpdateMovementRow(movement)
	wr, err := r.conn.WriteOneParameterizedContext(ctx, gorqlite.ParameterizedStatement{
		Query: `UPDATE movements SET occurred_on = ?, merchant = ?, merchant_key = ?, amount = ?, category = ?, method = ?, note = ?, edited = ?, updated_at = ?
		        WHERE id = ?`,
		Arguments: []any{row.OccurredOn, row.Merchant, merchantKey(row.Merchant), row.Amount, row.Category, row.Method, row.Note, row.Edited, now(), row.ID},
	})
	if err != nil {
		return domain.Movement{}, domainError(err, fmt.Sprintf("updating movement %d", row.ID))
	}
	if wr.RowsAffected == 0 {
		return domain.Movement{}, domain.ErrNotFound
	}

	qr, err := r.conn.QueryOneParameterizedContext(ctx, gorqlite.ParameterizedStatement{
		Query:     "SELECT id, ref, occurred_on, merchant, amount, category, method, note, source, edited, deleted_at FROM movements WHERE id = ?",
		Arguments: []any{row.ID},
	})
	if err != nil {
		return domain.Movement{}, domainError(err, fmt.Sprintf("selecting updated movement %d", row.ID))
	}
	if !qr.Next() {
		return domain.Movement{}, domain.ErrNotFound
	}
	var saved updateMovementRow
	if err := qr.Scan(&saved.ID, &saved.Ref, &saved.OccurredOn, &saved.Merchant, &saved.Amount, &saved.Category, &saved.Method, &saved.Note, &saved.Source, &saved.Edited, &saved.DeletedAt); err != nil {
		return domain.Movement{}, fmt.Errorf("reading movement %d: %w: %w", row.ID, domain.ErrInternal, err)
	}
	return saved.toDomain(), nil
}
