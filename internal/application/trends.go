package application

import (
	"context"
	"fmt"
	"sort"

	"github.com/drew138/abacus/internal/domain"
)

func (a *Application) Trends(ctx context.Context) (domain.Trends, error) {
	categories, err := a.repo.List(ctx)
	if err != nil {
		return domain.Trends{}, fmt.Errorf("listing categories: %w", err)
	}
	months, err := a.repo.Months(ctx)
	if err != nil {
		return domain.Trends{}, fmt.Errorf("listing months: %w", err)
	}
	totals, err := a.repo.Totals(ctx)
	if err != nil {
		return domain.Trends{}, fmt.Errorf("totalling months: %w", err)
	}
	return buildTrends(categories, months, totals), nil
}

func buildTrends(categories []domain.Category, months []domain.MonthSummary, totals []domain.CategoryTotal) domain.Trends {
	index := map[string]int{}
	for i, m := range months {
		index[m.Month] = i
	}
	parents := map[string]string{}
	for _, c := range categories {
		parents[c.Name] = c.Parent
	}

	tops := map[string]*domain.TrendRow{}
	subs := map[string]map[string]*domain.TrendRow{}
	var order []string
	for _, t := range totals {
		i, ok := index[t.Month]
		if !ok || t.Category == domain.CategorySavings || t.Total == 0 {
			continue
		}
		parent := parents[t.Category]
		top := domain.Category{Name: t.Category, Parent: parent}.Top()
		row, ok := tops[top]
		if !ok {
			row = &domain.TrendRow{Name: top, Values: make([]int64, len(months))}
			tops[top] = row
			subs[top] = map[string]*domain.TrendRow{}
			order = append(order, top)
		}
		row.Values[i] += t.Total
		if parent == "" {
			continue
		}
		sub, ok := subs[top][t.Category]
		if !ok {
			sub = &domain.TrendRow{Name: t.Category, Values: make([]int64, len(months))}
			subs[top][t.Category] = sub
		}
		sub.Values[i] += t.Total
	}

	last := len(months) - 1
	rows := make([]domain.TrendRow, 0, len(order))
	for _, name := range order {
		row := *tops[name]
		for _, sub := range subs[name] {
			row.Subcategories = append(row.Subcategories, *sub)
		}
		sort.Slice(row.Subcategories, func(i, j int) bool {
			return trendsLess(row.Subcategories[i], row.Subcategories[j], last, 0, 0)
		})
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return trendsLess(rows[i], rows[j], last, domain.CategoryTier(rows[i].Name, ""), domain.CategoryTier(rows[j].Name, ""))
	})
	return domain.Trends{Months: months, Rows: rows}
}

func trendsLess(a, b domain.TrendRow, last, tierA, tierB int) bool {
	if tierA != tierB {
		return tierA < tierB
	}
	if last >= 0 && a.Values[last] != b.Values[last] {
		return a.Values[last] > b.Values[last]
	}
	return a.Name < b.Name
}
