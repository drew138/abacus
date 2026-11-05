package application

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) MonthReport(ctx context.Context, month string, filter domain.ExpenseFilter) (domain.MonthReport, error) {
	categories, err := a.repo.List(ctx)
	if err != nil {
		return domain.MonthReport{}, fmt.Errorf("listing categories: %w", err)
	}
	months, err := a.repo.Months(ctx)
	if err != nil {
		return domain.MonthReport{}, fmt.Errorf("listing months: %w", err)
	}
	totals, err := a.repo.Totals(ctx)
	if err != nil {
		return domain.MonthReport{}, fmt.Errorf("totalling months: %w", err)
	}
	movements, err := a.repo.ListByMonth(ctx, month)
	if err != nil {
		return domain.MonthReport{}, fmt.Errorf("listing movements for %s: %w", month, err)
	}
	return buildMonthReport(month, filter, categories, months, totals, movements), nil
}

func buildMonthReport(month string, filter domain.ExpenseFilter, categories []domain.Category, months []domain.MonthSummary, totals []domain.CategoryTotal, movements []domain.Movement) domain.MonthReport {
	report := domain.MonthReport{Month: month, Summary: domain.MonthSummary{Month: month}, Filter: filter}
	for i, m := range months {
		if m.Month != month {
			continue
		}
		report.Summary = m
		if i > 0 {
			report.Previous, report.HasPrevious = months[i-1], true
		}
	}

	parents := map[string]string{}
	for _, c := range categories {
		parents[c.Name] = c.Parent
	}
	current, previous := map[string]int64{}, map[string]int64{}
	for _, t := range totals {
		switch {
		case t.Month == month:
			current[t.Category] = t.Total
		case report.HasPrevious && t.Month == report.Previous.Month:
			previous[t.Category] = t.Total
		}
	}

	report.Groups = monthReportGroups(parents, current, previous, report.HasPrevious)
	for _, m := range movements {
		if m.IsSavings() {
			report.Savings = append(report.Savings, m)
		}
	}
	report.Expenses = monthReportExpenses(parents, movements, filter)
	return report
}

func monthReportGroups(parents map[string]string, current, previous map[string]int64, hasPrevious bool) []domain.GroupTotal {
	byTop := map[string]*domain.GroupTotal{}
	var order []string
	for name, total := range current {
		if name == domain.CategorySavings || total == 0 {
			continue
		}
		top := domain.Category{Name: name, Parent: parents[name]}.Top()
		g, ok := byTop[top]
		if !ok {
			g = &domain.GroupTotal{CategoryAmount: domain.CategoryAmount{Name: top, HasPrevious: hasPrevious}}
			byTop[top] = g
			order = append(order, top)
		}
		g.Total += total
		if parents[name] != "" {
			g.Subcategories = append(g.Subcategories, domain.CategoryAmount{Name: name, Total: total, Previous: previous[name], HasPrevious: hasPrevious})
		}
	}
	for name, total := range previous {
		if g, ok := byTop[domain.Category{Name: name, Parent: parents[name]}.Top()]; ok && name != domain.CategorySavings {
			g.Previous += total
		}
	}

	groups := make([]domain.GroupTotal, 0, len(order))
	for _, top := range order {
		g := *byTop[top]
		sort.Slice(g.Subcategories, func(i, j int) bool {
			return monthReportLess(g.Subcategories[i].Name, g.Subcategories[j].Name, g.Subcategories[i].Total, g.Subcategories[j].Total, 0, 0)
		})
		groups = append(groups, g)
	}
	sort.Slice(groups, func(i, j int) bool {
		return monthReportLess(groups[i].Name, groups[j].Name, groups[i].Total, groups[j].Total, domain.CategoryTier(groups[i].Name, ""), domain.CategoryTier(groups[j].Name, ""))
	})
	return groups
}

func monthReportExpenses(parents map[string]string, movements []domain.Movement, filter domain.ExpenseFilter) []domain.ExpenseGroup {
	search := strings.ToLower(strings.TrimSpace(filter.Search))
	byCategory := map[string]*domain.ExpenseGroup{}
	var order []string
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
			g = &domain.ExpenseGroup{Category: m.Category, Parent: parent}
			byCategory[m.Category] = g
			order = append(order, m.Category)
		}
		g.Total += m.Amount
		g.Movements = append(g.Movements, m)
	}

	groups := make([]domain.ExpenseGroup, 0, len(order))
	for _, name := range order {
		g := *byCategory[name]
		sort.SliceStable(g.Movements, func(i, j int) bool {
			if !g.Movements[i].Date.Equal(g.Movements[j].Date) {
				return g.Movements[i].Date.Before(g.Movements[j].Date)
			}
			return g.Movements[i].Amount > g.Movements[j].Amount
		})
		groups = append(groups, g)
	}
	sort.Slice(groups, func(i, j int) bool {
		return monthReportLess(groups[i].Category, groups[j].Category, groups[i].Total, groups[j].Total,
			domain.CategoryTier(groups[i].Category, groups[i].Parent), domain.CategoryTier(groups[j].Category, groups[j].Parent))
	})
	return groups
}

func monthReportLess(nameA, nameB string, totalA, totalB int64, tierA, tierB int) bool {
	if tierA != tierB {
		return tierA < tierB
	}
	if totalA != totalB {
		return totalA > totalB
	}
	return nameA < nameB
}
