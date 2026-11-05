package sqlite

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) List(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT name, parent FROM categories ORDER BY name")
	if err != nil {
		return nil, domainError(err, "selecting categories")
	}
	defer rows.Close()

	categories := []domain.Category{}
	for rows.Next() {
		var row categoryRow
		if err := rows.Scan(&row.Name, &row.Parent); err != nil {
			return nil, domainError(err, "reading a category")
		}
		categories = append(categories, row.toDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading categories")
	}
	return categories, nil
}
