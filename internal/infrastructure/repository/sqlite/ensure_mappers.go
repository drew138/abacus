package sqlite

import (
	"database/sql"

	"github.com/drew138/abacus/internal/domain"
)

func mapEnsureCategoryRow(c domain.Category) ensureCategoryRow {
	return ensureCategoryRow{Name: c.Name, Parent: sql.NullString{String: c.Parent, Valid: c.Parent != ""}}
}
