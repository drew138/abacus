package rest

type getTrendsPage struct {
	Title   string
	Sidebar getTrendsSidebar
	Trends  getTrendsView
}

type getTrendsSidebar struct {
	Trends bool
	Months []getTrendsSidebarItem
}

type getTrendsSidebarItem struct {
	Key      string
	Name     string
	Badge    int
	Amount   string
	Latest   bool
	Selected bool
}

type getTrendsDelta struct {
	Text  string
	Class string
}

type getTrendsView struct {
	Empty     bool
	Subtitle  string
	Saved     string
	MonthsN   string
	Average   string
	LastName  string
	LastSpent string
	LastDelta getTrendsDelta
	Only      string
	Chart     getTrendsChart
	Legend    []getTrendsLegend
	Months    []string
	Savings   getTrendsTableRow
	Rows      []getTrendsTableRow
	Totals    getTrendsTableRow
}

type getTrendsLegend struct {
	Name  string
	Color string
	Off   bool
}

type getTrendsTableRow struct {
	Name  string
	Color string
	Cells []string
	Delta getTrendsDelta
	Subs  []getTrendsTableRow
}

type getTrendsChart struct {
	Width   int
	Height  int
	Grid    []getTrendsGridLine
	Columns []getTrendsColumn
}

type getTrendsGridLine struct {
	Y     string
	Label string
}

type getTrendsColumn struct {
	X        string
	LabelX   string
	TopY     string
	Total    string
	Month    string
	Segments []getTrendsSegment
}

type getTrendsSegment struct {
	Key   string
	Color string
	X     string
	Y     string
	W     string
	H     string
	Month string
	Value string
	Share string
	Delta getTrendsDelta
	Vs    string
}

type getTrendsRow struct {
	Name          string
	Values        []int64
	Subcategories []getTrendsRow
}
