package rest

type getRowView struct {
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
