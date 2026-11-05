package sqlite

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Totals(ctx context.Context) ([]domain.CategoryTotal, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT substr(occurred_on, 1, 7) AS month, category, sum(amount)
		 FROM movements WHERE deleted_at IS NULL
		 GROUP BY month, category ORDER BY month, category`)
	if err != nil {
		return nil, domainError(err, "selecting totals")
	}
	defer rows.Close()

	totals := []domain.CategoryTotal{}
	for rows.Next() {
		var row totalRow
		if err := rows.Scan(&row.Month, &row.Category, &row.Total); err != nil {
			return nil, domainError(err, "reading a total")
		}
		totals = append(totals, row.toDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading totals")
	}
	return totals, nil
}
