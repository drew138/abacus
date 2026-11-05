package rqlite

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) List(ctx context.Context) ([]domain.Category, error) {
	qr, err := r.conn.QueryOneContext(ctx, "SELECT name, parent FROM categories ORDER BY name")
	if err != nil {
		return nil, domainError(err, "selecting categories")
	}
	categories := []domain.Category{}
	for qr.Next() {
		var row listCategoryRow
		if err := qr.Scan(&row.Name, &row.Parent); err != nil {
			return nil, fmt.Errorf("reading a category: %w: %w", domain.ErrInternal, err)
		}
		categories = append(categories, row.toDomain())
	}
	return categories, nil
}
