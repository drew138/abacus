package postgres

import "context"

func (r *Repository) ListCategoryTotals(ctx context.Context) (map[string]map[string]int64, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT to_char(occurred_on, 'YYYY-MM') AS month, category, sum(amount)
		 FROM movements WHERE deleted_at IS NULL
		 GROUP BY month, category`)
	if err != nil {
		return nil, domainError(err, "selecting totals")
	}
	defer rows.Close()

	totals := map[string]map[string]int64{}
	for rows.Next() {
		var month, category string
		var total int64
		if err := rows.Scan(&month, &category, &total); err != nil {
			return nil, domainError(err, "reading a total")
		}
		if totals[month] == nil {
			totals[month] = map[string]int64{}
		}
		totals[month][category] = total
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading totals")
	}
	return totals, nil
}
