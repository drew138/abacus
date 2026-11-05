package rqlite

import (
	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

type listCategoryRow struct {
	Name   string
	Parent gorqlite.NullString
}

func (r listCategoryRow) toDomain() domain.Category {
	return domain.Category{Name: r.Name, Parent: r.Parent.String}
}
