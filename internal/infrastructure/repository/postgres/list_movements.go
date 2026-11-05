package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) ListMovements(ctx context.Context, month string) ([]domain.Movement, error) {
	start, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, fmt.Errorf("parsing month: %w: %w", domain.ErrInvalid, err)
	}
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, ref, occurred_on, merchant, amount, category, method, note, source, edited, deleted_at FROM movements WHERE deleted_at IS NULL AND occurred_on >= $1 AND occurred_on < $2 ORDER BY occurred_on, id",
		start.Format(dateLayout), start.AddDate(0, 1, 0).Format(dateLayout))
	if err != nil {
		return nil, domainError(err, "selecting movements of a month")
	}
	defer rows.Close()

	movements := []domain.Movement{}
	for rows.Next() {
		var row listMovementsRow
		if err := rows.Scan(&row.ID, &row.Ref, &row.OccurredOn, &row.Merchant, &row.Amount, &row.Category, &row.Method, &row.Note, &row.Source, &row.Edited, &row.DeletedAt); err != nil {
			return nil, domainError(err, "reading a movement")
		}
		movements = append(movements, row.toDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading movements")
	}
	return movements, nil
}
