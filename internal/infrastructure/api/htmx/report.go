package htmx

import (
	"sort"
	"strings"

	"github.com/drew138/abacus/internal/domain"
)

const (
	categoryOther = "Other"
	groupFood     = "Food and eating out"
)

type expenseFilter struct {
	Category string
	Search   string
}

type categoryAmount struct {
	Name        string
	Total       int64
	Previous    int64
	HasPrevious bool
}

type groupTotal struct {
	categoryAmount
	Subcategories []categoryAmount
}

type expenseGroup struct {
	Category  string
	Parent    string
	Total     int64
	Movements []domain.Movement
}

type monthReport struct {
	Month       string
	Summary     domain.MonthSummary
	Previous    domain.MonthSummary
	HasPrevious bool
	Groups      []groupTotal
	Savings     []domain.Movement
	Expenses    []expenseGroup
	Filter      expenseFilter
}

type trendRow struct {
	Name          string
	Values        []int64
	Subcategories []trendRow
}

type trends struct {
	Months []domain.MonthSummary
	Rows   []trendRow
}

func categoryTier(name, parent string) int {
	switch {
	case name == domain.CategorySavings:
		return 0
	case name == groupFood || parent == groupFood:
		return 1
	case name == categoryOther:
		return 3
	}
	return 2
}

func ordered(nameA, nameB string, totalA, totalB int64, tierA, tierB int) bool {
	if tierA != tierB {
		return tierA < tierB
	}
	if totalA != totalB {
		return totalA > totalB
	}
	return nameA < nameB
}

func parentsOf(categories []domain.Category) map[string]string {
	parents := make(map[string]string, len(categories))
	for _, c := range categories {
		parents[c.Name] = c.Parent
	}
	return parents
}

func buildMonthReport(month string, filter expenseFilter, categories []domain.Category, months []domain.MonthSummary, totals map[string]map[string]int64, movements []domain.Movement) monthReport {
	report := monthReport{Month: month, Summary: domain.MonthSummary{Month: month}, Filter: filter}
	for i, m := range months {
		if m.Month != month {
			continue
		}
		report.Summary = m
		if i > 0 {
			report.Previous, report.HasPrevious = months[i-1], true
		}
	}

	parents := parentsOf(categories)
	previous := map[string]int64{}
	if report.HasPrevious {
		previous = totals[report.Previous.Month]
	}
	report.Groups = buildGroups(parents, totals[month], previous, report.HasPrevious)
	for _, m := range movements {
		if m.IsSavings() {
			report.Savings = append(report.Savings, m)
		}
	}
	report.Expenses = buildExpenses(parents, movements, filter)
	return report
}

func buildGroups(parents map[string]string, current, previous map[string]int64, hasPrevious bool) []groupTotal {
	byTop := map[string]*groupTotal{}
	for name, total := range current {
		if name == domain.CategorySavings || total == 0 {
			continue
		}
		top := domain.Category{Name: name, Parent: parents[name]}.Top()
		g, ok := byTop[top]
		if !ok {
			g = &groupTotal{categoryAmount: categoryAmount{Name: top, HasPrevious: hasPrevious}}
			byTop[top] = g
		}
		g.Total += total
		if parents[name] != "" {
			g.Subcategories = append(g.Subcategories, categoryAmount{Name: name, Total: total, Previous: previous[name], HasPrevious: hasPrevious})
		}
	}
	for name, total := range previous {
		if g, ok := byTop[domain.Category{Name: name, Parent: parents[name]}.Top()]; ok && name != domain.CategorySavings {
			g.Previous += total
		}
	}

	groups := make([]groupTotal, 0, len(byTop))
	for _, g := range byTop {
		sort.Slice(g.Subcategories, func(i, j int) bool {
			return ordered(g.Subcategories[i].Name, g.Subcategories[j].Name, g.Subcategories[i].Total, g.Subcategories[j].Total, 0, 0)
		})
		groups = append(groups, *g)
	}
	sort.Slice(groups, func(i, j int) bool {
		return ordered(groups[i].Name, groups[j].Name, groups[i].Total, groups[j].Total, categoryTier(groups[i].Name, ""), categoryTier(groups[j].Name, ""))
	})
	return groups
}

func buildExpenses(parents map[string]string, movements []domain.Movement, filter expenseFilter) []expenseGroup {
	search := strings.ToLower(strings.TrimSpace(filter.Search))
	byCategory := map[string]*expenseGroup{}
	for _, m := range movements {
		parent := parents[m.Category]
		if filter.Category != "" && m.Category != filter.Category && parent != filter.Category {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(m.Merchant), search) {
			continue
		}
		g, ok := byCategory[m.Category]
		if !ok {
			g = &expenseGroup{Category: m.Category, Parent: parent}
			byCategory[m.Category] = g
		}
		g.Total += m.Amount
		g.Movements = append(g.Movements, m)
	}

	groups := make([]expenseGroup, 0, len(byCategory))
	for _, g := range byCategory {
		sort.SliceStable(g.Movements, func(i, j int) bool {
			if !g.Movements[i].Date.Equal(g.Movements[j].Date) {
				return g.Movements[i].Date.Before(g.Movements[j].Date)
			}
			return g.Movements[i].Amount > g.Movements[j].Amount
		})
		groups = append(groups, *g)
	}
	sort.Slice(groups, func(i, j int) bool {
		return ordered(groups[i].Category, groups[j].Category, groups[i].Total, groups[j].Total,
			categoryTier(groups[i].Category, groups[i].Parent), categoryTier(groups[j].Category, groups[j].Parent))
	})
	return groups
}

func buildTrends(categories []domain.Category, months []domain.MonthSummary, totals map[string]map[string]int64) trends {
	parents := parentsOf(categories)
	tops := map[string]*trendRow{}
	subs := map[string]map[string]*trendRow{}
	for i, m := range months {
		for name, total := range totals[m.Month] {
			if name == domain.CategorySavings || total == 0 {
				continue
			}
			parent := parents[name]
			top := domain.Category{Name: name, Parent: parent}.Top()
			row, ok := tops[top]
			if !ok {
				row = &trendRow{Name: top, Values: make([]int64, len(months))}
				tops[top] = row
				subs[top] = map[string]*trendRow{}
			}
			row.Values[i] += total
			if parent == "" {
				continue
			}
			sub, ok := subs[top][name]
			if !ok {
				sub = &trendRow{Name: name, Values: make([]int64, len(months))}
				subs[top][name] = sub
			}
			sub.Values[i] += total
		}
	}

	last := len(months) - 1
	latest := func(r trendRow) int64 {
		if last < 0 {
			return 0
		}
		return r.Values[last]
	}
	rows := make([]trendRow, 0, len(tops))
	for name, row := range tops {
		for _, sub := range subs[name] {
			row.Subcategories = append(row.Subcategories, *sub)
		}
		sort.Slice(row.Subcategories, func(i, j int) bool {
			a, b := row.Subcategories[i], row.Subcategories[j]
			return ordered(a.Name, b.Name, latest(a), latest(b), 0, 0)
		})
		rows = append(rows, *row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return ordered(rows[i].Name, rows[j].Name, latest(rows[i]), latest(rows[j]), categoryTier(rows[i].Name, ""), categoryTier(rows[j].Name, ""))
	})
	return trends{Months: months, Rows: rows}
}
