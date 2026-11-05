package rest

import (
	"sort"

	"github.com/drew138/abacus/internal/domain"
)

type menuView struct {
	ID       int64
	Merchant string
	Month    string
	Options  []menuOptionView
}

type menuOptionView struct {
	Name    string
	Color   string
	Current bool
}

func mapMenuView(m domain.Movement, categories []domain.Category) menuView {
	view := menuView{ID: m.ID, Merchant: m.Merchant, Month: m.Month()}
	sorted := append([]domain.Category(nil), categories...)
	sort.Slice(sorted, func(i, j int) bool {
		ti, tj := categoryTier(sorted[i].Name, sorted[i].Parent), categoryTier(sorted[j].Name, sorted[j].Parent)
		if ti != tj {
			return ti < tj
		}
		return sorted[i].Name < sorted[j].Name
	})
	for _, c := range sorted {
		if c.Name == domain.CategorySavings {
			continue
		}
		view.Options = append(view.Options, menuOptionView{Name: c.Name, Color: categoryColor(c.Name), Current: c.Name == m.Category})
	}
	return view
}
