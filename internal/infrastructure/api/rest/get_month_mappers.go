package rest

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"unicode"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/pesos"
)

func mapGetMonthPage(month string, categories []domain.Category, months []domain.MonthSummary, totals map[string]map[string]int64, movements []domain.Movement) getMonthPage {
	summary := domain.MonthSummary{Month: month}
	var previous domain.MonthSummary
	hasPrevious := false
	for i, m := range months {
		if m.Month == month {
			summary = m
			if i > 0 {
				previous, hasPrevious = months[i-1], true
			}
		}
	}
	parents := make(map[string]string, len(categories))
	for _, c := range categories {
		parents[c.Name] = c.Parent
	}

	view := getMonthView{
		Key:         month,
		Title:       monthName(month),
		Spent:       pesos.Format(summary.Spent),
		Purchases:   summary.Movements,
		HasPrevious: hasPrevious,
		Expenses:    mapGetMonthExpenses(month, parents, movements),
	}
	previousTotals := map[string]int64{}
	if hasPrevious {
		view.Delta = mapGetMonthDelta(summary.Spent, previous.Spent)
		view.PreviousName = monthShort(previous.Month)
		previousTotals = totals[previous.Month]
	}

	groups := mapGetMonthGroups(parents, totals[month], previousTotals)
	var largest int64 = 1
	for _, g := range groups {
		largest = max(largest, g.Total)
	}
	for _, g := range groups {
		view.Ranks = append(view.Ranks, mapGetMonthRank(g.Name, g.Total, g.Previous, hasPrevious, largest, summary.Spent, false, g.Name))
		for _, sub := range g.Subcategories {
			view.Ranks = append(view.Ranks, mapGetMonthRank(sub.Name, sub.Total, sub.Previous, hasPrevious, largest, summary.Spent, true, g.Name))
		}
	}
	view.Donut = mapGetMonthDonut(groups, summary)

	var saved int64
	for _, m := range movements {
		if !m.IsSavings() {
			continue
		}
		saved += m.Amount
		view.Savings = append(view.Savings, getMonthSaving{Label: fmt.Sprintf("%d %s · %s", m.Date.Day(), m.Date.Format("Jan"), m.Merchant), Amount: pesos.Format(m.Amount)})
	}
	view.SavingsTotal = pesos.Format(saved)

	return getMonthPage{Title: view.Title, Sidebar: mapGetMonthSidebar(months, month), Month: view}
}

func mapGetMonthSidebar(months []domain.MonthSummary, selected string) getMonthSidebar {
	sidebar := getMonthSidebar{}
	for i := len(months) - 1; i >= 0; i-- {
		m := months[i]
		sidebar.Months = append(sidebar.Months, getMonthSidebarItem{
			Key:      m.Month,
			Name:     monthName(m.Month),
			Badge:    m.Movements,
			Amount:   pesos.Short(m.Spent),
			Latest:   i == len(months)-1,
			Selected: m.Month == selected,
		})
	}
	return sidebar
}

func mapGetMonthDelta(current, previous int64) getMonthDelta {
	if previous == 0 {
		return getMonthDelta{}
	}
	p := float64(current-previous) / float64(previous) * 100
	switch {
	case math.Abs(p) < 1:
		return getMonthDelta{Text: "±0%"}
	case p > 0:
		return getMonthDelta{Text: fmt.Sprintf("▲ %.0f%%", p), Class: "up"}
	}
	return getMonthDelta{Text: fmt.Sprintf("▼ %.0f%%", -p), Class: "down"}
}

func mapGetMonthGroups(parents map[string]string, current, previous map[string]int64) []getMonthGroup {
	byTop := map[string]*getMonthGroup{}
	for name, total := range current {
		if name == domain.CategorySavings || total == 0 {
			continue
		}
		top := domain.Category{Name: name, Parent: parents[name]}.Top()
		g, ok := byTop[top]
		if !ok {
			g = &getMonthGroup{Name: top}
			byTop[top] = g
		}
		g.Total += total
		if parents[name] != "" {
			g.Subcategories = append(g.Subcategories, getMonthAmount{Name: name, Total: total, Previous: previous[name]})
		}
	}
	for name, total := range previous {
		if g, ok := byTop[domain.Category{Name: name, Parent: parents[name]}.Top()]; ok && name != domain.CategorySavings {
			g.Previous += total
		}
	}

	groups := make([]getMonthGroup, 0, len(byTop))
	for _, g := range byTop {
		sort.Slice(g.Subcategories, func(i, j int) bool {
			a, b := g.Subcategories[i], g.Subcategories[j]
			if a.Total != b.Total {
				return a.Total > b.Total
			}
			return a.Name < b.Name
		})
		groups = append(groups, *g)
	}
	sort.Slice(groups, func(i, j int) bool {
		a, b := groups[i], groups[j]
		if ta, tb := categoryTier(a.Name, ""), categoryTier(b.Name, ""); ta != tb {
			return ta < tb
		}
		if a.Total != b.Total {
			return a.Total > b.Total
		}
		return a.Name < b.Name
	})
	return groups
}

func mapGetMonthRank(name string, total, previous int64, hasPrevious bool, largest, spent int64, sub bool, group string) getMonthRank {
	rank := getMonthRank{
		Name:   name,
		Color:  categoryColor(name),
		Amount: pesos.Format(total),
		Width:  strconv.FormatFloat(float64(total)*100/float64(largest), 'f', 1, 64) + "%",
		Share:  pesos.Percent(total, spent) + " of the month",
		Sub:    sub,
		Filter: name,
		Hover:  group,
	}
	if hasPrevious {
		rank.Delta = mapGetMonthDelta(total, previous)
	}
	return rank
}

func mapGetMonthDonut(groups []getMonthGroup, summary domain.MonthSummary) getMonthDonut {
	const outer, inner, cx, cy = 84.0, 60.0, 96.0, 96.0
	donut := getMonthDonut{Total: pesos.Format(summary.Spent), Count: summary.Movements}
	var total int64
	for _, g := range groups {
		total += g.Total
	}
	if total == 0 {
		return donut
	}
	gap := 0.012
	if len(groups) == 1 {
		gap = 0.0001
	}
	point := func(angle, radius float64) string {
		return fmt.Sprintf("%.2f,%.2f", cx+radius*math.Cos(angle), cy+radius*math.Sin(angle))
	}
	start := -math.Pi / 2
	for _, g := range groups {
		fraction := float64(g.Total) / float64(total)
		end := start + fraction*2*math.Pi - gap
		large := 0
		if fraction > 0.5 {
			large = 1
		}
		mid := (start + end) / 2
		donut.Slices = append(donut.Slices, getMonthSlice{
			Key:    g.Name,
			Name:   g.Name,
			Color:  categoryColor(g.Name),
			Path:   fmt.Sprintf("M%s A%.0f,%.0f 0 %d 1 %s L%s A%.0f,%.0f 0 %d 0 %sZ", point(start, outer), outer, outer, large, point(end, outer), point(end, inner), inner, inner, large, point(start, inner)),
			DX:     fmt.Sprintf("%.1fpx", math.Cos(mid)*6),
			DY:     fmt.Sprintf("%.1fpx", math.Sin(mid)*6),
			Short:  pesos.Short(g.Total),
			Share:  pesos.Percent(g.Total, total),
			Filter: g.Name,
			Amount: pesos.Format(g.Total),
		})
		start += fraction * 2 * math.Pi
	}
	return donut
}

func mapGetMonthExpenses(month string, parents map[string]string, movements []domain.Movement) getMonthExpenses {
	expenses := getMonthExpenses{Month: month}
	byCategory := map[string]*getMonthExpenseGroup{}
	totals := map[string]int64{}
	sorted := append([]domain.Movement(nil), movements...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if !sorted[i].Date.Equal(sorted[j].Date) {
			return sorted[i].Date.Before(sorted[j].Date)
		}
		return sorted[i].Amount > sorted[j].Amount
	})
	for _, m := range sorted {
		g, ok := byCategory[m.Category]
		if !ok {
			g = &getMonthExpenseGroup{Category: m.Category, Parent: parents[m.Category], Color: categoryColor(m.Category), IsSavings: m.IsSavings()}
			byCategory[m.Category] = g
		}
		totals[m.Category] += m.Amount
		g.Rows = append(g.Rows, mapGetMonthRow(m))
		g.Count++
		expenses.Count++
	}
	for _, g := range byCategory {
		g.Total = pesos.Format(totals[g.Category])
		expenses.Groups = append(expenses.Groups, *g)
	}
	sort.Slice(expenses.Groups, func(i, j int) bool {
		a, b := expenses.Groups[i], expenses.Groups[j]
		if ta, tb := categoryTier(a.Category, a.Parent), categoryTier(b.Category, b.Parent); ta != tb {
			return ta < tb
		}
		if totals[a.Category] != totals[b.Category] {
			return totals[a.Category] > totals[b.Category]
		}
		return a.Category < b.Category
	})
	if expenses.Count == 0 {
		expenses.Empty = "No movements in this month yet."
	}
	return expenses
}

func mapGetMonthRow(m domain.Movement) getMonthRow {
	initial := "?"
	for _, r := range m.Merchant {
		initial = string(unicode.ToUpper(r))
		break
	}
	return getMonthRow{
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
