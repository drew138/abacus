package rest

type getExpensesQuery struct {
	Category string
	Search   string
}

type getExpensesView struct {
	Month  string
	Count  int
	Filter string
	Search string
	Groups []getExpensesGroup
	Empty  string
}

type getExpensesGroup struct {
	Category  string
	Parent    string
	Color     string
	Count     int
	Total     string
	IsSavings bool
	Rows      []getExpensesRow
}

type getExpensesRow struct {
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
