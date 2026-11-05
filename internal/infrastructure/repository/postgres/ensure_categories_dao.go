package postgres

import "database/sql"

type ensureCategoriesRow struct {
	Name   string
	Parent sql.NullString
}
