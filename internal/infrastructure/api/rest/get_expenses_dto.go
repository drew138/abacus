package rest

type getExpensesQuery struct {
	Category string `json:"category"`
	Search   string `json:"search"`
}

type getExpensesView struct {
	Month  string             `json:"month"`
	Count  int                `json:"count"`
	Filter string             `json:"filter"`
	Search string             `json:"search"`
	Groups []getExpensesGroup `json:"groups"`
	Empty  string             `json:"empty"`
}

type getExpensesGroup struct {
	Category  string           `json:"category"`
	Parent    string           `json:"parent"`
	Color     string           `json:"color"`
	Count     int              `json:"count"`
	Total     string           `json:"total"`
	IsSavings bool             `json:"is_savings"`
	Rows      []getExpensesRow `json:"rows"`
}

type getExpensesRow struct {
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
