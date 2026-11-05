package postgres

import (
	"database/sql"

	"github.com/drew138/abacus/internal/domain"
)

type listCategoriesRow struct {
	Name   string
	Parent sql.NullString
}

func (r listCategoriesRow) toDomain() domain.Category {
	return domain.Category{Name: r.Name, Parent: r.Parent.String}
}
