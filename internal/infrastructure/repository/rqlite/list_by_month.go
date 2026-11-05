package rqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) ListByMonth(ctx context.Context, month string) ([]domain.Movement, error) {
	start, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, fmt.Errorf("parsing month: %w: %w", domain.ErrInvalid, err)
	}
	qr, err := r.conn.QueryOneParameterizedContext(ctx, gorqlite.ParameterizedStatement{
		Query: `SELECT id, ref, occurred_on, merchant, amount, category, method, note, source, edited, deleted_at
		        FROM movements WHERE deleted_at IS NULL AND occurred_on >= ? AND occurred_on < ? ORDER BY occurred_on, id`,
		Arguments: []any{start.Format(dateLayout), start.AddDate(0, 1, 0).Format(dateLayout)},
	})
	if err != nil {
		return nil, domainError(err, "selecting movements of a month")
	}
	movements := []domain.Movement{}
	for qr.Next() {
		var row listByMonthMovementRow
		if err := qr.Scan(&row.ID, &row.Ref, &row.OccurredOn, &row.Merchant, &row.Amount, &row.Category, &row.Method, &row.Note, &row.Source, &row.Edited, &row.DeletedAt); err != nil {
			return nil, fmt.Errorf("reading a movement: %w: %w", domain.ErrInternal, err)
		}
		movements = append(movements, row.toDomain())
	}
	return movements, nil
}
