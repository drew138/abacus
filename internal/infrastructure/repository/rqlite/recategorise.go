package rqlite

import (
	"context"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Recategorise(ctx context.Context, merchant, category string) (int, error) {
	key, at := merchantKey(merchant), now()
	results, err := r.conn.WriteParameterizedContext(ctx, []gorqlite.ParameterizedStatement{
		{
			Query: `INSERT INTO merchant_rules (merchant_key, category, updated_at) VALUES (?, ?, ?)
			        ON CONFLICT (merchant_key) DO UPDATE SET category = excluded.category, updated_at = excluded.updated_at`,
			Arguments: []any{key, category, at},
		},
		{
			Query:     "UPDATE movements SET category = ?, updated_at = ? WHERE merchant_key = ? AND category <> ? AND category <> ?",
			Arguments: []any{category, at, key, category, domain.CategorySavings},
		},
	})
	if err != nil {
		return 0, domainError(err, "recategorising a merchant")
	}
	return int(results[1].RowsAffected), nil
}
