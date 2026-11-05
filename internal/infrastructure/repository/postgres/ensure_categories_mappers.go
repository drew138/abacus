package postgres

import (
	"database/sql"

	"github.com/drew138/abacus/internal/domain"
)

func mapEnsureCategoriesRow(c domain.Category) ensureCategoriesRow {
	return ensureCategoriesRow{Name: c.Name, Parent: sql.NullString{String: c.Parent, Valid: c.Parent != ""}}
}
