package rest

type getMonthPage struct {
	Title   string          `json:"title"`
	Sidebar getMonthSidebar `json:"sidebar"`
	Month   getMonthView    `json:"month"`
}

type getMonthSidebar struct {
	Trends bool                  `json:"trends"`
	Months []getMonthSidebarItem `json:"months"`
}

type getMonthSidebarItem struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Badge    int    `json:"badge"`
	Amount   string `json:"amount"`
	Latest   bool   `json:"latest"`
	Selected bool   `json:"selected"`
}

type getMonthDelta struct {
	Text  string `json:"text"`
	Class string `json:"class"`
}

type getMonthView struct {
	Key          string           `json:"key"`
	Title        string           `json:"title"`
	Spent        string           `json:"spent"`
	Purchases    int              `json:"purchases"`
	HasPrevious  bool             `json:"has_previous"`
	Delta        getMonthDelta    `json:"delta"`
	PreviousName string           `json:"previous_name"`
	Donut        getMonthDonut    `json:"donut"`
	Ranks        []getMonthRank   `json:"ranks"`
	SavingsTotal string           `json:"savings_total"`
	Savings      []getMonthSaving `json:"savings"`
	Expenses     getMonthExpenses `json:"expenses"`
}

type getMonthDonut struct {
	Slices []getMonthSlice `json:"slices"`
	Total  string          `json:"total"`
	Count  int             `json:"count"`
}

type getMonthSlice struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	Path   string `json:"path"`
	DX     string `json:"dx"`
	DY     string `json:"dy"`
	Short  string `json:"short"`
	Share  string `json:"share"`
	Filter string `json:"filter"`
	Amount string `json:"amount"`
}

type getMonthRank struct {
	Name   string        `json:"name"`
	Color  string        `json:"color"`
	Amount string        `json:"amount"`
	Width  string        `json:"width"`
	Share  string        `json:"share"`
	Delta  getMonthDelta `json:"delta"`
	Sub    bool          `json:"sub"`
	Filter string        `json:"filter"`
	Hover  string        `json:"hover"`
}

type getMonthSaving struct {
	Label  string `json:"label"`
	Amount string `json:"amount"`
}

type getMonthExpenses struct {
	Month  string                 `json:"month"`
	Count  int                    `json:"count"`
	Filter string                 `json:"filter"`
	Search string                 `json:"search"`
	Groups []getMonthExpenseGroup `json:"groups"`
	Empty  string                 `json:"empty"`
}

type getMonthExpenseGroup struct {
	Category  string        `json:"category"`
	Parent    string        `json:"parent"`
	Color     string        `json:"color"`
	Count     int           `json:"count"`
	Total     string        `json:"total"`
	IsSavings bool          `json:"is_savings"`
	Rows      []getMonthRow `json:"rows"`
}

type getMonthRow struct {
	ID        int64  `json:"id"`
	Initial   string `json:"initial"`
	Color     string `json:"color"`
	Merchant  string `json:"merchant"`
	Day       string `json:"day"`
	Category  string `json:"category"`
	IsSavings bool   `json:"is_savings"`
	Method    string `json:"method"`
	Manual    bool   `json:"manual"`
	Edited    bool   `json:"edited"`
	Note      string `json:"note"`
	Amount    string `json:"amount"`
	Confirm   bool   `json:"confirm"`
	Month     string `json:"month"`
}

type getMonthGroup struct {
	Name          string           `json:"name"`
	Total         int64            `json:"total"`
	Previous      int64            `json:"previous"`
	Subcategories []getMonthAmount `json:"subcategories"`
}

type getMonthAmount struct {
	Name     string `json:"name"`
	Total    int64  `json:"total"`
	Previous int64  `json:"previous"`
}
