package rest

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/pesos"
)

func mapGetExpensesView(month string, query getExpensesQuery, categories []domain.Category, movements []domain.Movement) getExpensesView {
	parents := make(map[string]string, len(categories))
	for _, c := range categories {
		parents[c.Name] = c.Parent
	}
	view := getExpensesView{Month: month, Filter: query.Category, Search: query.Search}
	search := strings.ToLower(strings.TrimSpace(query.Search))
	sorted := append([]domain.Movement(nil), movements...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if !sorted[i].Date.Equal(sorted[j].Date) {
			return sorted[i].Date.Before(sorted[j].Date)
		}
		return sorted[i].Amount > sorted[j].Amount
	})

	byCategory := map[string]*getExpensesGroup{}
	totals := map[string]int64{}
	for _, m := range sorted {
		parent := parents[m.Category]
		if query.Category != "" && m.Category != query.Category && parent != query.Category {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(m.Merchant), search) {
			continue
		}
		g, ok := byCategory[m.Category]
		if !ok {
			g = &getExpensesGroup{Category: m.Category, Parent: parent, Color: categoryColor(m.Category), IsSavings: m.IsSavings()}
			byCategory[m.Category] = g
		}
		totals[m.Category] += m.Amount
		g.Rows = append(g.Rows, mapGetExpensesRow(m))
		g.Count++
		view.Count++
	}
	for _, g := range byCategory {
		g.Total = pesos.Format(totals[g.Category])
		view.Groups = append(view.Groups, *g)
	}
	sort.Slice(view.Groups, func(i, j int) bool {
		a, b := view.Groups[i], view.Groups[j]
		if ta, tb := categoryTier(a.Category, a.Parent), categoryTier(b.Category, b.Parent); ta != tb {
			return ta < tb
		}
		if totals[a.Category] != totals[b.Category] {
			return totals[a.Category] > totals[b.Category]
		}
		return a.Category < b.Category
	})
	if view.Count == 0 {
		view.Empty = "No movements in this month yet."
		if query.Category != "" || query.Search != "" {
			view.Empty = "No movements match. Clear the search or the category filter."
		}
	}
	return view
}

func mapGetExpensesRow(m domain.Movement) getExpensesRow {
	initial := "?"
	for _, r := range m.Merchant {
		initial = string(unicode.ToUpper(r))
		break
	}
	return getExpensesRow{
		ID:        m.ID,
		Initial:   initial,
		Color:     categoryColor(m.Category),
		Merchant:  m.Merchant,
		Day:       fmt.Sprintf("%d %s", m.Date.Day(), m.Date.Format("Jan")),
		Category:  m.Category,
		IsSavings: m.IsSavings(),
		Method:    methodLabels[m.Method],
		Manual:    m.Source == domain.SourceManual,
		Edited:    m.Edited,
		Note:      m.Note,
		Amount:    pesos.Format(m.Amount),
		Month:     m.Month(),
	}
}
