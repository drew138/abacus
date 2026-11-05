package rqlite

import (
	"context"
	"fmt"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Find(ctx context.Context, id int64) (domain.Movement, error) {
	qr, err := r.conn.QueryOneParameterizedContext(ctx, gorqlite.ParameterizedStatement{
		Query:     "SELECT id, ref, occurred_on, merchant, amount, category, method, note, source, edited, deleted_at FROM movements WHERE id = ?",
		Arguments: []any{id},
	})
	if err != nil {
		return domain.Movement{}, domainError(err, fmt.Sprintf("selecting movement %d", id))
	}
	if !qr.Next() {
		return domain.Movement{}, domain.ErrNotFound
	}
	var row findMovementRow
	if err := qr.Scan(&row.ID, &row.Ref, &row.OccurredOn, &row.Merchant, &row.Amount, &row.Category, &row.Method, &row.Note, &row.Source, &row.Edited, &row.DeletedAt); err != nil {
		return domain.Movement{}, fmt.Errorf("reading movement %d: %w: %w", id, domain.ErrInternal, err)
	}
	return row.toDomain(), nil
}
