package rqlite

import (
	"context"
	"fmt"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Months(ctx context.Context) ([]domain.MonthSummary, error) {
	qr, err := r.conn.QueryOneParameterizedContext(ctx, gorqlite.ParameterizedStatement{
		Query: `SELECT substr(occurred_on, 1, 7) AS month,
		               coalesce(sum(CASE WHEN category = ? THEN 0 ELSE amount END), 0),
		               coalesce(sum(CASE WHEN category = ? THEN amount ELSE 0 END), 0),
		               coalesce(sum(CASE WHEN category = ? THEN 0 ELSE 1 END), 0)
		        FROM movements WHERE deleted_at IS NULL
		        GROUP BY month ORDER BY month`,
		Arguments: []any{domain.CategorySavings, domain.CategorySavings, domain.CategorySavings},
	})
	if err != nil {
		return nil, domainError(err, "selecting months")
	}
	months := []domain.MonthSummary{}
	for qr.Next() {
		var row monthsRow
		if err := qr.Scan(&row.Month, &row.Spent, &row.Saved, &row.Movements); err != nil {
			return nil, fmt.Errorf("reading a month: %w: %w", domain.ErrInternal, err)
		}
		months = append(months, row.toDomain())
	}
	return months, nil
}
