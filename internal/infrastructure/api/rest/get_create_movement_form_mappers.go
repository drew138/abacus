package rest

import (
	"sort"
	"strings"
	"time"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/pesos"
)

func mapGetCreateMovementFormView(m domain.Movement, categories []domain.Category, today time.Time, errs []string) getCreateMovementFormView {
	view := getCreateMovementFormView{
		ID:            m.ID,
		Title:         "Add movement",
		Merchant:      m.Merchant,
		Date:          m.Date.Format("2006-01-02"),
		MaxDate:       today.Format("2006-01-02"),
		Note:          m.Note,
		FromStatement: m.Source == domain.SourceStatement,
		Errors:        errs,
	}
	if m.Amount > 0 {
		view.Amount = strings.TrimPrefix(pesos.Format(m.Amount), "$")
	}
	for _, method := range []domain.Method{domain.MethodCredit, domain.MethodDebit} {
		view.Methods = append(view.Methods, getCreateMovementFormOption{Value: string(method), Label: methodLabels[method], Selected: method == m.Method})
	}
	sorted := append([]domain.Category(nil), categories...)
	sort.Slice(sorted, func(i, j int) bool {
		ti, tj := categoryTier(sorted[i].Name, sorted[i].Parent), categoryTier(sorted[j].Name, sorted[j].Parent)
		if ti != tj {
			return ti < tj
		}
		return sorted[i].Name < sorted[j].Name
	})
	for _, c := range sorted {
		option := getCreateMovementFormOption{Value: c.Name, Label: c.Name, Color: categoryColor(c.Name), Selected: c.Name == m.Category}
		last := len(view.Categories) - 1
		if last >= 0 && view.Categories[last].Label == c.Parent {
			view.Categories[last].Options = append(view.Categories[last].Options, option)
			continue
		}
		view.Categories = append(view.Categories, getCreateMovementFormOptionGroup{Label: c.Parent, Options: []getCreateMovementFormOption{option}})
	}
	return view
}
