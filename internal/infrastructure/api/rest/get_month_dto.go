package rest

type getMonthPage struct {
	Title   string
	Sidebar getMonthSidebar
	Month   getMonthView
}

type getMonthSidebar struct {
	Trends bool
	Months []getMonthSidebarItem
}

type getMonthSidebarItem struct {
	Key      string
	Name     string
	Badge    int
	Amount   string
	Latest   bool
	Selected bool
}

type getMonthDelta struct {
	Text  string
	Class string
}

type getMonthView struct {
	Key          string
	Title        string
	Spent        string
	Purchases    int
	HasPrevious  bool
	Delta        getMonthDelta
	PreviousName string
	Donut        getMonthDonut
	Ranks        []getMonthRank
	SavingsTotal string
	Savings      []getMonthSaving
	Expenses     getMonthExpenses
}

type getMonthDonut struct {
	Slices []getMonthSlice
	Total  string
	Count  int
}

type getMonthSlice struct {
	Key    string
	Name   string
	Color  string
	Path   string
	DX     string
	DY     string
	Short  string
	Share  string
	Filter string
	Amount string
}

type getMonthRank struct {
	Name   string
	Color  string
	Amount string
	Width  string
	Share  string
	Delta  getMonthDelta
	Sub    bool
	Filter string
	Hover  string
}

type getMonthSaving struct {
	Label  string
	Amount string
}

type getMonthExpenses struct {
	Month  string
	Count  int
	Filter string
	Search string
	Groups []getMonthExpenseGroup
	Empty  string
}

type getMonthExpenseGroup struct {
	Category  string
	Parent    string
	Color     string
	Count     int
	Total     string
	IsSavings bool
	Rows      []getMonthRow
}

type getMonthRow struct {
	ID        int64
	Initial   string
	Color     string
	Merchant  string
	Day       string
	Category  string
	IsSavings bool
	Method    string
	Manual    bool
	Edited    bool
	Note      string
	Amount    string
	Confirm   bool
	Month     string
}

type getMonthGroup struct {
	Name          string
	Total         int64
	Previous      int64
	Subcategories []getMonthAmount
}

type getMonthAmount struct {
	Name     string
	Total    int64
	Previous int64
}
