package sqlite

import (
	"database/sql"

	"github.com/drew138/abacus/internal/domain"
)

type listCategoryRow struct {
	Name   string
	Parent sql.NullString
}

func (r listCategoryRow) toDomain() domain.Category {
	return domain.Category{Name: r.Name, Parent: r.Parent.String}
}
