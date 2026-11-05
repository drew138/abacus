package sqlite

import "database/sql"

type ensureCategoryRow struct {
	Name   string
	Parent sql.NullString
}
