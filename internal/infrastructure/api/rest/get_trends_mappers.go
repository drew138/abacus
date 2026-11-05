package rest

import (
	"fmt"
	"math"
	"sort"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/pesos"
)

func mapGetTrendsPage(categories []domain.Category, months []domain.MonthSummary, totals map[string]map[string]int64, only string) getTrendsPage {
	sidebar := getTrendsSidebar{Trends: true}
	for i := len(months) - 1; i >= 0; i-- {
		m := months[i]
		sidebar.Months = append(sidebar.Months, getTrendsSidebarItem{Key: m.Month, Name: monthName(m.Month), Badge: m.Movements, Amount: pesos.Short(m.Spent), Latest: i == len(months)-1})
	}
	return getTrendsPage{Title: "Month to month", Sidebar: sidebar, Trends: mapGetTrendsView(categories, months, totals, only)}
}

func mapGetTrendsView(categories []domain.Category, months []domain.MonthSummary, totals map[string]map[string]int64, only string) getTrendsView {
	n := len(months)
	if n == 0 {
		return getTrendsView{Empty: true}
	}
	rows := mapGetTrendsRows(categories, months, totals)
	view := getTrendsView{Only: only}
	var saved, spent int64
	for _, m := range months {
		saved += m.Saved
		spent += m.Spent
		view.Months = append(view.Months, monthShort(m.Month))
	}
	last := months[n-1]
	average := spent / int64(n)
	unit := "months"
	if n == 1 {
		unit = "month"
	}
	view.Subtitle = fmt.Sprintf("%s to %s · %d %s", monthName(months[0].Month), monthName(last.Month), n, unit)
	view.Saved = pesos.Format(saved)
	view.MonthsN = fmt.Sprintf("across %d %s", n, unit)
	view.Average = pesos.Format(average)
	view.LastName = monthShort(last.Month)
	view.LastSpent = pesos.Format(last.Spent)
	view.LastDelta = mapGetTrendsDelta(last.Spent, average, false)

	view.Savings = getTrendsTableRow{Name: "Moved to savings", Color: categoryColor(domain.CategorySavings)}
	view.Totals = getTrendsTableRow{Name: "Total spent"}
	for _, m := range months {
		view.Savings.Cells = append(view.Savings.Cells, mapGetTrendsCell(m.Saved))
		view.Totals.Cells = append(view.Totals.Cells, pesos.Short(m.Spent))
	}
	if n > 1 {
		view.Savings.Delta = mapGetTrendsDelta(last.Saved, months[n-2].Saved, true)
		view.Totals.Delta = mapGetTrendsDelta(last.Spent, months[n-2].Spent, false)
	}

	for _, row := range rows {
		view.Legend = append(view.Legend, getTrendsLegend{Name: row.Name, Color: categoryColor(row.Name), Off: only != "" && only != row.Name})
		table := mapGetTrendsTableRow(row)
		for _, sub := range row.Subcategories {
			table.Subs = append(table.Subs, mapGetTrendsTableRow(sub))
		}
		view.Rows = append(view.Rows, table)
	}
	view.Chart = mapGetTrendsChart(months, rows, only)
	return view
}

func mapGetTrendsRows(categories []domain.Category, months []domain.MonthSummary, totals map[string]map[string]int64) []getTrendsRow {
	parents := make(map[string]string, len(categories))
	for _, c := range categories {
		parents[c.Name] = c.Parent
	}
	tops := map[string]*getTrendsRow{}
	subs := map[string]map[string]*getTrendsRow{}
	for i, m := range months {
		for name, total := range totals[m.Month] {
			if name == domain.CategorySavings || total == 0 {
				continue
			}
			parent := parents[name]
			top := domain.Category{Name: name, Parent: parent}.Top()
			row, ok := tops[top]
			if !ok {
				row = &getTrendsRow{Name: top, Values: make([]int64, len(months))}
				tops[top] = row
				subs[top] = map[string]*getTrendsRow{}
			}
			row.Values[i] += total
			if parent == "" {
				continue
			}
			sub, ok := subs[top][name]
			if !ok {
				sub = &getTrendsRow{Name: name, Values: make([]int64, len(months))}
				subs[top][name] = sub
			}
			sub.Values[i] += total
		}
	}

	last := len(months) - 1
	less := func(a, b getTrendsRow, tierA, tierB int) bool {
		if tierA != tierB {
			return tierA < tierB
		}
		if a.Values[last] != b.Values[last] {
			return a.Values[last] > b.Values[last]
		}
		return a.Name < b.Name
	}
	rows := make([]getTrendsRow, 0, len(tops))
	for name, row := range tops {
		for _, sub := range subs[name] {
			row.Subcategories = append(row.Subcategories, *sub)
		}
		sort.Slice(row.Subcategories, func(i, j int) bool { return less(row.Subcategories[i], row.Subcategories[j], 0, 0) })
		rows = append(rows, *row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return less(rows[i], rows[j], categoryTier(rows[i].Name, ""), categoryTier(rows[j].Name, ""))
	})
	return rows
}

func mapGetTrendsDelta(current, previous int64, higherIsGood bool) getTrendsDelta {
	if previous == 0 {
		return getTrendsDelta{}
	}
	p := float64(current-previous) / float64(previous) * 100
	if math.Abs(p) < 1 {
		return getTrendsDelta{Text: "±0%"}
	}
	up := p > 0
	class := "down"
	if up != higherIsGood {
		class = "up"
	}
	arrow := "▼"
	if up {
		arrow = "▲"
	}
	return getTrendsDelta{Text: fmt.Sprintf("%s %.0f%%", arrow, math.Abs(p)), Class: class}
}

func mapGetTrendsCell(v int64) string {
	if v == 0 {
		return "–"
	}
	return pesos.Short(v)
}

func mapGetTrendsTableRow(row getTrendsRow) getTrendsTableRow {
	table := getTrendsTableRow{Name: row.Name, Color: categoryColor(row.Name)}
	for _, v := range row.Values {
		table.Cells = append(table.Cells, mapGetTrendsCell(v))
	}
	if n := len(row.Values); n > 1 {
		table.Delta = mapGetTrendsDelta(row.Values[n-1], row.Values[n-2], false)
	}
	return table
}

func mapGetTrendsChart(months []domain.MonthSummary, rows []getTrendsRow, only string) getTrendsChart {
	const width, height, left, bottom, top = 640, 250, 56.0, 26.0, 16.0
	n := len(months)
	chart := getTrendsChart{Width: width, Height: height}
	monthTotals := make([]int64, n)
	shown := make([]int64, n)
	for _, row := range rows {
		for i, v := range row.Values {
			monthTotals[i] += v
			if only == "" || only == row.Name {
				shown[i] += v
			}
		}
	}
	var highest int64 = 1
	for _, v := range shown {
		highest = max(highest, v)
	}
	step := math.Pow(10, math.Floor(math.Log10(float64(highest)/4)))
	nice := step * 10
	for _, f := range []float64{1, 2, 2.5, 5, 10} {
		if float64(highest)/(f*step) <= 5 {
			nice = f * step
			break
		}
	}
	ceiling := math.Ceil(float64(highest)/nice) * nice
	y := func(v float64) float64 { return top + (height-bottom-top)*(1-v/ceiling) }
	for g := 0.0; g <= ceiling+1; g += nice {
		chart.Grid = append(chart.Grid, getTrendsGridLine{Y: fmt.Sprintf("%.1f", y(g)), Label: pesos.Short(int64(g))})
	}

	columnWidth := (width - left - 10) / float64(n)
	barWidth := math.Min(columnWidth*0.6, 64)
	for i, m := range months {
		x := left + float64(i)*columnWidth + (columnWidth-barWidth)/2
		column := getTrendsColumn{X: fmt.Sprintf("%.1f", x), LabelX: fmt.Sprintf("%.1f", x+barWidth/2), Total: pesos.Short(shown[i]), Month: monthShort(m.Month)}
		var stack float64
		for _, row := range rows {
			v := row.Values[i]
			if v == 0 || (only != "" && only != row.Name) {
				continue
			}
			segment := getTrendsSegment{
				Key:   row.Name,
				Color: categoryColor(row.Name),
				X:     column.X,
				Y:     fmt.Sprintf("%.1f", y(stack+float64(v))),
				W:     fmt.Sprintf("%.1f", barWidth),
				H:     fmt.Sprintf("%.1f", math.Max(0, y(stack)-y(stack+float64(v))-1)),
				Month: monthName(m.Month),
				Value: pesos.Format(v),
			}
			if only == "" {
				segment.Share = pesos.Percent(v, monthTotals[i]) + " of the month"
			}
			if i > 0 && row.Values[i-1] > 0 {
				segment.Delta = mapGetTrendsDelta(v, row.Values[i-1], false)
				segment.Vs = "vs " + monthShort(months[i-1].Month)
			}
			column.Segments = append(column.Segments, segment)
			stack += float64(v)
		}
		column.TopY = fmt.Sprintf("%.1f", y(stack)-6)
		chart.Columns = append(chart.Columns, column)
	}
	return chart
}
