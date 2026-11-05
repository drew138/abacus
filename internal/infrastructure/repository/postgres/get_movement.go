package postgres

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) GetMovement(ctx context.Context, id int64) (domain.Movement, error) {
	var row getMovementRow
	err := r.db.QueryRowContext(ctx, "SELECT id, ref, occurred_on, merchant, amount, category, method, note, source, edited, deleted_at FROM movements WHERE id = $1", id).
		Scan(&row.ID, &row.Ref, &row.OccurredOn, &row.Merchant, &row.Amount, &row.Category, &row.Method, &row.Note, &row.Source, &row.Edited, &row.DeletedAt)
	if err != nil {
		return domain.Movement{}, domainError(err, fmt.Sprintf("selecting movement %d", id))
	}
	return row.toDomain(), nil
}
