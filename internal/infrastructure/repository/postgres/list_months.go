package postgres

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) ListMonths(ctx context.Context) ([]domain.MonthSummary, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT to_char(occurred_on, 'YYYY-MM') AS month,
		        coalesce(sum(amount) FILTER (WHERE category <> $1), 0),
		        coalesce(sum(amount) FILTER (WHERE category = $1), 0),
		        count(*) FILTER (WHERE category <> $1)
		 FROM movements WHERE deleted_at IS NULL
		 GROUP BY month ORDER BY month`, domain.CategorySavings)
	if err != nil {
		return nil, domainError(err, "selecting months")
	}
	defer rows.Close()

	months := []domain.MonthSummary{}
	for rows.Next() {
		var row listMonthsRow
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
