package rqlite

import "github.com/drew138/abacus/internal/domain"

func mapEnsureCategoryRow(c domain.Category) ensureCategoryRow {
	var parent any
	if c.Parent != "" {
		parent = c.Parent
	}
	return ensureCategoryRow{Name: c.Name, Parent: parent}
}
