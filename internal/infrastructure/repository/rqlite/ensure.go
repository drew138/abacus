package rqlite

import (
	"context"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Ensure(ctx context.Context, categories []domain.Category) error {
	if len(categories) == 0 {
		return nil
	}
	at := now()
	statements := make([]gorqlite.ParameterizedStatement, 0, len(categories))
	for _, c := range categories {
		row := mapEnsureCategoryRow(c)
		statements = append(statements, gorqlite.ParameterizedStatement{
			Query:     "INSERT INTO categories (name, parent, created_at) VALUES (?, ?, ?) ON CONFLICT (name) DO NOTHING",
			Arguments: []any{row.Name, row.Parent, at},
		})
	}
	if _, err := r.conn.WriteParameterizedContext(ctx, statements); err != nil {
		return domainError(err, "inserting categories")
	}
	return nil
}
