package sqlite

import (
	"database/sql"

	"github.com/drew138/abacus/internal/domain"
)

func mapCategoryRow(c domain.Category) categoryRow {
	return categoryRow{Name: c.Name, Parent: sql.NullString{String: c.Parent, Valid: c.Parent != ""}}
}
