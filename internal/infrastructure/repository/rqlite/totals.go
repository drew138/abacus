package rqlite

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Totals(ctx context.Context) (map[string]map[string]int64, error) {
	qr, err := r.conn.QueryOneContext(ctx,
		`SELECT substr(occurred_on, 1, 7) AS month, category, sum(amount)
		 FROM movements WHERE deleted_at IS NULL
		 GROUP BY month, category`)
	if err != nil {
		return nil, domainError(err, "selecting totals")
	}
	totals := map[string]map[string]int64{}
	for qr.Next() {
		var month, category string
		var total int64
		if err := qr.Scan(&month, &category, &total); err != nil {
			return nil, fmt.Errorf("reading a total: %w: %w", domain.ErrInternal, err)
		}
		if totals[month] == nil {
			totals[month] = map[string]int64{}
		}
		totals[month][category] = total
	}
	return totals, nil
}
