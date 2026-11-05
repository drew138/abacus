package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"

	"github.com/drew138/abacus/internal/domain"
)

var knownColors = map[string]string{
	"Restaurants and cafés":             "#e0763a",
	"Food delivery":                     "#f0a35e",
	"Groceries":                         "#c9a227",
	"Snacks and candy":                  "#b7791f",
	"Transfers to people":               "#d9534f",
	"Rides and transport":               "#7c5cdb",
	"Gym and fitness":                   "#1f9d74",
	"Phone, internet and subscriptions": "#2b6fd6",
	"Online shopping":                   "#d14d8b",
	"Pharmacy and health":               "#0e9fb0",
	"Gaming":                            "#5865c9",
	"Entertainment":                     "#e87ba4",
	"Travel":                            "#3d8fb0",
	categoryOther:                       "#8a909c",
	domain.CategorySavings:              "#15803d",
	groupFood:                           "#e0763a",
}

var spareColors = []string{"#a0622d", "#b5179e", "#6f8f00", "#3d5a99", "#d98e04", "#008300", "#e34948", "#0096b4"}

var methodLabels = []optionView{
	{Value: string(domain.MethodCredit), Label: "Credit Card"},
	{Value: string(domain.MethodDebit), Label: "Bank Transfer"},
}

func categoryColor(name string) string {
	if c, ok := knownColors[name]; ok {
		return c
	}
	var h uint32
	for _, r := range name {
		h = h*31 + uint32(r)
	}
	return spareColors[h%uint32(len(spareColors))]
}

func methodLabel(m domain.Method) string {
	for _, o := range methodLabels {
		if o.Value == string(m) {
			return o.Label
		}
	}
	return string(m)
}

func pesos(v int64) string {
	sign := ""
	if v < 0 {
		sign, v = "-", -v
	}
	digits := strconv.FormatInt(v, 10)
	var b strings.Builder
	for i, d := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(d)
	}
	return sign + "$" + b.String()
}

func shortPesos(v int64) string {
	switch {
	case v >= 1_000_000:
		s := strconv.FormatFloat(float64(v)/1_000_000, 'f', 1, 64)
		return "$" + strings.Replace(strings.TrimSuffix(s, ".0"), ".", ",", 1) + " M"
	case v >= 1_000:
		return "$" + strconv.FormatInt(int64(math.Round(float64(v)/1_000)), 10) + " k"
	}
	return "$" + strconv.FormatInt(v, 10)
}

func percent(part, whole int64) string {
	if whole == 0 {
		return "0%"
	}
	return strings.Replace(strconv.FormatFloat(float64(part)*100/float64(whole), 'f', 1, 64), ".", ",", 1) + "%"
}

func mapDelta(current, previous int64, higherIsGood bool) deltaView {
	if previous == 0 {
		return deltaView{}
	}
	p := float64(current-previous) / float64(previous) * 100
	if math.Abs(p) < 1 {
		return deltaView{Text: "±0%"}
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
	return deltaView{Text: fmt.Sprintf("%s %.0f%%", arrow, math.Abs(p)), Class: class}
}

func monthTime(key string) time.Time {
	t, _ := time.Parse("2006-01", key)
	return t
}

func monthName(key string) string {
	return monthTime(key).Format("January 2006")
}

func monthShort(key string) string {
	return monthTime(key).Format("Jan")
}

func mapSidebarView(months []domain.MonthSummary, selected string) sidebarView {
	view := sidebarView{Trends: selected == ""}
	for i := len(months) - 1; i >= 0; i-- {
		m := months[i]
		view.Months = append(view.Months, sidebarMonth{
			Key:      m.Month,
			Name:     monthName(m.Month),
			Badge:    m.Movements,
			Amount:   shortPesos(m.Spent),
			Latest:   i == len(months)-1,
			Selected: m.Month == selected,
		})
	}
	return view
}

func mapMonthView(report monthReport) monthView {
	view := monthView{
		Key:          report.Month,
		Title:        monthName(report.Month),
		Spent:        pesos(report.Summary.Spent),
		Purchases:    report.Summary.Movements,
		HasPrevious:  report.HasPrevious,
		SavingsTotal: pesos(0),
		Expenses:     mapExpensesView(report),
	}
	if report.HasPrevious {
		view.Delta = mapDelta(report.Summary.Spent, report.Previous.Spent, false)
		view.PreviousName = monthShort(report.Previous.Month)
	}

	var maxTotal int64 = 1
	for _, g := range report.Groups {
		maxTotal = max(maxTotal, g.Total)
	}
	for _, g := range report.Groups {
		filter := g.Name
		view.Ranks = append(view.Ranks, mapRankView(g.categoryAmount, maxTotal, report.Summary.Spent, false, filter, g.Name))
		for _, sub := range g.Subcategories {
			view.Ranks = append(view.Ranks, mapRankView(sub, maxTotal, report.Summary.Spent, true, sub.Name, g.Name))
		}
	}
	view.Donut = mapDonutView(report.Groups, report.Summary)

	var saved int64
	for _, m := range report.Savings {
		saved += m.Amount
		view.Savings = append(view.Savings, savingView{Label: fmt.Sprintf("%d %s · %s", m.Date.Day(), m.Date.Format("Jan"), m.Merchant), Amount: pesos(m.Amount)})
	}
	view.SavingsTotal = pesos(saved)
	return view
}

func mapRankView(a categoryAmount, maxTotal, spent int64, sub bool, filter, hover string) rankView {
	view := rankView{
		Name:   a.Name,
		Color:  categoryColor(a.Name),
		Amount: pesos(a.Total),
		Width:  strconv.FormatFloat(float64(a.Total)*100/float64(maxTotal), 'f', 1, 64) + "%",
		Share:  percent(a.Total, spent) + " of the month",
		Sub:    sub,
		Filter: filter,
		Hover:  hover,
	}
	if a.HasPrevious {
		view.Delta = mapDelta(a.Total, a.Previous, false)
	}
	return view
}

func mapDonutView(groups []groupTotal, summary domain.MonthSummary) donutView {
	const outer, inner, cx, cy = 84.0, 60.0, 96.0, 96.0
	view := donutView{Total: pesos(summary.Spent), Count: summary.Movements}
	var total int64
	for _, g := range groups {
		total += g.Total
	}
	if total == 0 {
		return view
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
		view.Slices = append(view.Slices, sliceView{
			Key:    g.Name,
			Name:   g.Name,
			Color:  categoryColor(g.Name),
			Path:   fmt.Sprintf("M%s A%.0f,%.0f 0 %d 1 %s L%s A%.0f,%.0f 0 %d 0 %sZ", point(start, outer), outer, outer, large, point(end, outer), point(end, inner), inner, inner, large, point(start, inner)),
			DX:     fmt.Sprintf("%.1fpx", math.Cos(mid)*6),
			DY:     fmt.Sprintf("%.1fpx", math.Sin(mid)*6),
			Short:  shortPesos(g.Total),
			Share:  percent(g.Total, total),
			Filter: g.Name,
			Amount: pesos(g.Total),
		})
		start += fraction * 2 * math.Pi
	}
	return view
}

func mapExpensesView(report monthReport) expensesView {
	view := expensesView{Month: report.Month, Filter: report.Filter.Category, Search: report.Filter.Search}
	for _, g := range report.Expenses {
		group := expenseGroupView{Category: g.Category, Parent: g.Parent, Color: categoryColor(g.Category), Count: len(g.Movements), Total: pesos(g.Total), IsSavings: g.Category == domain.CategorySavings}
		for _, m := range g.Movements {
			group.Rows = append(group.Rows, mapRowView(m, false))
		}
		view.Count += len(g.Movements)
		view.Groups = append(view.Groups, group)
	}
	if view.Count == 0 {
		view.Empty = "No movements in this month yet."
		if report.Filter.Category != "" || report.Filter.Search != "" {
			view.Empty = "No movements match. Clear the search or the category filter."
		}
	}
	return view
}

func mapRowView(m domain.Movement, confirm bool) rowView {
	initial := "?"
	for _, r := range m.Merchant {
		initial = string(unicode.ToUpper(r))
		break
	}
	return rowView{
		ID:        m.ID,
		Initial:   initial,
		Color:     categoryColor(m.Category),
		Merchant:  m.Merchant,
		Day:       fmt.Sprintf("%d %s", m.Date.Day(), m.Date.Format("Jan")),
		Category:  m.Category,
		IsSavings: m.IsSavings(),
		Method:    methodLabel(m.Method),
		Manual:    m.Source == domain.SourceManual,
		Edited:    m.Edited,
		Note:      m.Note,
		Amount:    pesos(m.Amount),
		Confirm:   confirm,
		Month:     m.Month(),
	}
}

func mapTrendsView(trends trends, only string) trendsView {
	n := len(trends.Months)
	if n == 0 {
		return trendsView{Empty: true}
	}
	view := trendsView{Only: only}
	var saved, spent int64
	for _, m := range trends.Months {
		saved += m.Saved
		spent += m.Spent
		view.Months = append(view.Months, monthShort(m.Month))
	}
	last := trends.Months[n-1]
	average := spent / int64(n)
	monthsWord := "months"
	if n == 1 {
		monthsWord = "month"
	}
	view.Subtitle = fmt.Sprintf("%s to %s · %d %s", monthName(trends.Months[0].Month), monthName(last.Month), n, monthsWord)
	view.Saved = pesos(saved)
	view.MonthsN = fmt.Sprintf("across %d %s", n, monthsWord)
	view.Average = pesos(average)
	view.LastName = monthShort(last.Month)
	view.LastSpent = pesos(last.Spent)
	view.LastDelta = mapDelta(last.Spent, average, false)

	savings := tableRowView{Name: "Moved to savings", Color: categoryColor(domain.CategorySavings)}
	totals := tableRowView{Name: "Total spent"}
	for _, m := range trends.Months {
		savings.Cells = append(savings.Cells, cell(m.Saved))
		totals.Cells = append(totals.Cells, shortPesos(m.Spent))
	}
	if n > 1 {
		savings.Delta = mapDelta(last.Saved, trends.Months[n-2].Saved, true)
		totals.Delta = mapDelta(last.Spent, trends.Months[n-2].Spent, false)
	}
	view.Savings, view.Totals = savings, totals

	for _, row := range trends.Rows {
		view.Legend = append(view.Legend, legendView{Name: row.Name, Color: categoryColor(row.Name), Off: only != "" && only != row.Name})
		table := mapTableRowView(row)
		for _, sub := range row.Subcategories {
			table.Subs = append(table.Subs, mapTableRowView(sub))
		}
		view.Rows = append(view.Rows, table)
	}
	view.Chart = mapBarsView(trends, only)
	return view
}

func cell(v int64) string {
	if v == 0 {
		return "–"
	}
	return shortPesos(v)
}

func mapTableRowView(row trendRow) tableRowView {
	view := tableRowView{Name: row.Name, Color: categoryColor(row.Name)}
	for _, v := range row.Values {
		view.Cells = append(view.Cells, cell(v))
	}
	if n := len(row.Values); n > 1 {
		view.Delta = mapDelta(row.Values[n-1], row.Values[n-2], false)
	}
	return view
}

func mapBarsView(trends trends, only string) barsView {
	const width, height, left, bottom, top = 640, 250, 56.0, 26.0, 16.0
	n := len(trends.Months)
	view := barsView{Width: width, Height: height}
	monthTotals := make([]int64, n)
	shown := make([]int64, n)
	for _, row := range trends.Rows {
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
		view.Grid = append(view.Grid, gridView{Y: fmt.Sprintf("%.1f", y(g)), Label: shortPesos(int64(g))})
	}

	columnWidth := (width - left - 10) / float64(n)
	barWidth := math.Min(columnWidth*0.6, 64)
	for i, m := range trends.Months {
		x := left + float64(i)*columnWidth + (columnWidth-barWidth)/2
		column := columnView{X: fmt.Sprintf("%.1f", x), LabelX: fmt.Sprintf("%.1f", x+barWidth/2), Total: shortPesos(shown[i]), Month: monthShort(m.Month)}
		var stack float64
		for _, row := range trends.Rows {
			v := row.Values[i]
			if v == 0 || (only != "" && only != row.Name) {
				continue
			}
			segment := segmentView{
				Key:   row.Name,
				Color: categoryColor(row.Name),
				X:     column.X,
				Y:     fmt.Sprintf("%.1f", y(stack+float64(v))),
				W:     fmt.Sprintf("%.1f", barWidth),
				H:     fmt.Sprintf("%.1f", math.Max(0, y(stack)-y(stack+float64(v))-1)),
				Month: monthName(m.Month),
				Value: pesos(v),
			}
			if only == "" {
				segment.Share = percent(v, monthTotals[i]) + " of the month"
			}
			if i > 0 && row.Values[i-1] > 0 {
				segment.Delta = mapDelta(v, row.Values[i-1], false)
				segment.Vs = "vs " + monthShort(trends.Months[i-1].Month)
			}
			column.Segments = append(column.Segments, segment)
			stack += float64(v)
		}
		column.TopY = fmt.Sprintf("%.1f", y(stack)-6)
		view.Columns = append(view.Columns, column)
	}
	return view
}

func mapFormView(m domain.Movement, categories []domain.Category, today time.Time, errs []string) formView {
	view := formView{
		ID:            m.ID,
		Title:         "Add movement",
		Merchant:      m.Merchant,
		Date:          m.Date.Format("2006-01-02"),
		MaxDate:       today.Format("2006-01-02"),
		Method:        string(m.Method),
		Note:          m.Note,
		FromStatement: m.Source == domain.SourceStatement,
		Errors:        errs,
	}
	if m.ID != 0 {
		view.Title = "Edit movement"
	}
	if m.Amount > 0 {
		view.Amount = strings.TrimPrefix(pesos(m.Amount), "$")
	}
	for _, o := range methodLabels {
		view.Methods = append(view.Methods, optionView{Value: o.Value, Label: o.Label, Selected: o.Value == string(m.Method)})
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
		option := optionView{Value: c.Name, Label: c.Name, Color: categoryColor(c.Name), Selected: c.Name == m.Category}
		last := len(view.Categories) - 1
		if last >= 0 && view.Categories[last].Label == c.Parent {
			view.Categories[last].Options = append(view.Categories[last].Options, option)
			continue
		}
		view.Categories = append(view.Categories, optionGroupView{Label: c.Parent, Options: []optionView{option}})
	}
	return view
}

func mapMovementForm(values url.Values) MovementForm {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, values.Get("amount"))
	amount, _ := strconv.ParseInt(digits, 10, 64)
	return MovementForm{
		Merchant:        strings.TrimSpace(values.Get("merchant")),
		Amount:          amount,
		Date:            values.Get("date"),
		Category:        values.Get("category"),
		Method:          values.Get("method"),
		Note:            strings.TrimSpace(values.Get("note")),
		ApplyToMerchant: values.Get("apply_to_merchant") == "on",
	}
}

func mapFormErrors(err error) []string {
	var failed validator.ValidationErrors
	if errors.As(err, &failed) {
		seen := map[string]bool{}
		var out []string
		for _, fe := range failed {
			if msg, ok := formFieldMessages[fe.Field()]; ok && !seen[msg] {
				seen[msg] = true
				out = append(out, msg)
			}
		}
		return out
	}
	if errors.Is(err, domain.ErrUnknownCategory) {
		return []string{formFieldMessages["category"]}
	}
	return nil
}

func mapErrorPage(err error) (int, string) {
	switch {
	case errors.Is(err, errMalformedForm):
		return http.StatusBadRequest, "That form could not be read. Reload the page and try again."
	case errors.Is(err, errCrossSite):
		return http.StatusForbidden, "That request did not come from Abacus. Reload the page and try again."
	case errors.Is(err, errInvalidMonth):
		return http.StatusBadRequest, "That month does not exist. Pick one from the sidebar."
	case errors.Is(err, errInvalidID):
		return http.StatusBadRequest, "That movement does not exist. Reload the month."
	case errors.Is(err, domain.ErrUnknownCategory):
		return http.StatusBadRequest, "Pick a category from the list."
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "That movement no longer exists. Reload the month."
	case errors.Is(err, domain.ErrInvalid):
		return http.StatusBadRequest, "Abacus could not accept that. Check the values and try again."
	case errors.Is(err, domain.ErrUnavailable), errors.Is(err, context.DeadlineExceeded):
		return http.StatusServiceUnavailable, "Abacus is busy right now. Try again in a moment."
	}
	return http.StatusInternalServerError, "Something went wrong. Try again, and reload the page if it keeps happening."
}

var domainErrors = []struct {
	err    error
	status int
	code   string
	field  string
}{
	{domain.ErrUnknownCategory, http.StatusBadRequest, "unknown_category", "category"},
}

var errorKinds = []struct {
	err    error
	status int
	typ    string
	title  string
}{
	{domain.ErrNotFound, http.StatusNotFound, "not_found", "The resource does not exist."},
	{domain.ErrConflict, http.StatusConflict, "resource_conflict", "The request conflicts with existing data."},
	{domain.ErrInvalid, http.StatusBadRequest, "invalid_argument", "The request is invalid."},
	{domain.ErrUnavailable, http.StatusServiceUnavailable, "service_unavailable", "The service is temporarily unavailable."},
	{context.DeadlineExceeded, http.StatusServiceUnavailable, "service_unavailable", "The service is temporarily unavailable."},
}

var tagCodes = map[string]string{
	"required":         "value_required",
	"required_if":      "value_required",
	"required_unless":  "value_required",
	"required_with":    "value_required",
	"required_without": "value_required",
	"oneof":            "invalid_option",
	"enum":             "invalid_option",
	"datetime":         "invalid_date",
}

var lengthCodes = map[string]string{
	"max": "max_length",
	"lte": "max_length",
	"min": "min_length",
	"gte": "min_length",
	"len": "length_mismatch",
}

var numberCodes = map[string]string{
	"max": "less_than_or_equal",
	"lte": "less_than_or_equal",
	"lt":  "less_than",
	"min": "greater_than_or_equal",
	"gte": "greater_than_or_equal",
	"gt":  "greater_than",
}

func mapErrorResponse(err error) ErrorResponse {
	var wrongType *json.UnmarshalTypeError
	var failed validator.ValidationErrors
	switch {
	case errors.As(err, &wrongType):
		return mapInvalidParamResponse(wrongType.Field, jsonType(wrongType.Type))
	case errors.Is(err, errMalformedBody):
		return ErrorResponse{Type: "malformed_request", Title: "The request body is not valid JSON.", Status: http.StatusBadRequest}
	case errors.As(err, &failed):
		return ErrorResponse{Type: "invalid_argument", Title: "The request has invalid fields.", Status: http.StatusBadRequest, Errors: mapErrorDetails(failed)}
	case errors.Is(err, errBadImportToken):
		return ErrorResponse{Type: "unauthenticated", Title: "A valid import token is required.", Status: http.StatusUnauthorized}
	}

	var details []ErrorDetail
	status := 0
	for _, d := range domainErrors {
		if errors.Is(err, d.err) {
			if status == 0 {
				status = d.status
			}
			details = append(details, ErrorDetail{Code: d.code, Field: d.field})
		}
	}
	for _, k := range errorKinds {
		if errors.Is(err, k.err) {
			if status == 0 {
				status = k.status
			}
			return ErrorResponse{Type: k.typ, Title: k.title, Status: status, Errors: details}
		}
	}
	return ErrorResponse{Type: "internal_error", Title: "Something went wrong.", Status: http.StatusInternalServerError}
}

func mapInvalidParamResponse(field, expected string) ErrorResponse {
	return ErrorResponse{Type: "invalid_argument", Title: "The request has invalid fields.", Status: http.StatusBadRequest,
		Errors: []ErrorDetail{{Code: "invalid_type", Field: field, Param: expected}}}
}

func mapErrorDetails(failed validator.ValidationErrors) []ErrorDetail {
	details := make([]ErrorDetail, 0, len(failed))
	for _, fe := range failed {
		field := fe.Field()
		if _, path, ok := strings.Cut(fe.Namespace(), "."); ok {
			field = path
		}
		details = append(details, mapErrorDetail(fe, field))
	}
	return details
}

func mapErrorDetail(fe validator.FieldError, field string) ErrorDetail {
	switch fe.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		if code, ok := lengthCodes[fe.Tag()]; ok {
			return ErrorDetail{Code: code, Field: field, Param: fe.Param()}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		if code, ok := numberCodes[fe.Tag()]; ok {
			return ErrorDetail{Code: code, Field: field, Param: fe.Param()}
		}
	}
	if code, ok := tagCodes[fe.Tag()]; ok {
		return ErrorDetail{Code: code, Field: field}
	}
	return ErrorDetail{Code: "invalid_value", Field: field}
}
