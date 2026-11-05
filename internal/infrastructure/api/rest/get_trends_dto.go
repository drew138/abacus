package rest

type getTrendsPage struct {
	Title   string           `json:"title"`
	Sidebar getTrendsSidebar `json:"sidebar"`
	Trends  getTrendsView    `json:"trends"`
}

type getTrendsSidebar struct {
	Trends bool                   `json:"trends"`
	Months []getTrendsSidebarItem `json:"months"`
}

type getTrendsSidebarItem struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Badge    int    `json:"badge"`
	Amount   string `json:"amount"`
	Latest   bool   `json:"latest"`
	Selected bool   `json:"selected"`
}

type getTrendsDelta struct {
	Text  string `json:"text"`
	Class string `json:"class"`
}

type getTrendsView struct {
	Empty     bool                `json:"empty"`
	Subtitle  string              `json:"subtitle"`
	Saved     string              `json:"saved"`
	MonthsN   string              `json:"months_n"`
	Average   string              `json:"average"`
	LastName  string              `json:"last_name"`
	LastSpent string              `json:"last_spent"`
	LastDelta getTrendsDelta      `json:"last_delta"`
	Only      string              `json:"only"`
	Chart     getTrendsChart      `json:"chart"`
	Legend    []getTrendsLegend   `json:"legend"`
	Months    []string            `json:"months"`
	Savings   getTrendsTableRow   `json:"savings"`
	Rows      []getTrendsTableRow `json:"rows"`
	Totals    getTrendsTableRow   `json:"totals"`
}

type getTrendsLegend struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Off   bool   `json:"off"`
}

type getTrendsTableRow struct {
	Name  string              `json:"name"`
	Color string              `json:"color"`
	Cells []string            `json:"cells"`
	Delta getTrendsDelta      `json:"delta"`
	Subs  []getTrendsTableRow `json:"subs"`
}

type getTrendsChart struct {
	Width   int                 `json:"width"`
	Height  int                 `json:"height"`
	Grid    []getTrendsGridLine `json:"grid"`
	Columns []getTrendsColumn   `json:"columns"`
}

type getTrendsGridLine struct {
	Y     string `json:"y"`
	Label string `json:"label"`
}

type getTrendsColumn struct {
	X        string             `json:"x"`
	LabelX   string             `json:"label_x"`
	TopY     string             `json:"top_y"`
	Total    string             `json:"total"`
	Month    string             `json:"month"`
	Segments []getTrendsSegment `json:"segments"`
}

type getTrendsSegment struct {
	Key   string         `json:"key"`
	Color string         `json:"color"`
	X     string         `json:"x"`
	Y     string         `json:"y"`
	W     string         `json:"w"`
	H     string         `json:"h"`
	Month string         `json:"month"`
	Value string         `json:"value"`
	Share string         `json:"share"`
	Delta getTrendsDelta `json:"delta"`
	Vs    string         `json:"vs"`
}

type getTrendsRow struct {
	Name          string         `json:"name"`
	Values        []int64        `json:"values"`
	Subcategories []getTrendsRow `json:"subcategories"`
}
