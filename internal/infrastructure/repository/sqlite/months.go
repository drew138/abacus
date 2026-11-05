package sqlite

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Months(ctx context.Context) ([]domain.MonthSummary, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT substr(occurred_on, 1, 7) AS month,
		        coalesce(sum(CASE WHEN category = ? THEN 0 ELSE amount END), 0),
		        coalesce(sum(CASE WHEN category = ? THEN amount ELSE 0 END), 0),
		        coalesce(sum(CASE WHEN category = ? THEN 0 ELSE 1 END), 0)
		 FROM movements WHERE deleted_at IS NULL
		 GROUP BY month ORDER BY month`,
		domain.CategorySavings, domain.CategorySavings, domain.CategorySavings)
	if err != nil {
		return nil, domainError(err, "selecting months")
	}
	defer rows.Close()

	var months []domain.MonthSummary
	for rows.Next() {
		var row monthRow
		if err := rows.Scan(&row.Month, &row.Spent, &row.Saved, &row.Movements); err != nil {
			return nil, domainError(err, "reading a month")
		}
		months = append(months, row.toDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading months")
	}
	return months, nil
}
