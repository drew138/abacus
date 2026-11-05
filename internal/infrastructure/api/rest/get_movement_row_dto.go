package rest

type getMovementRowView struct {
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
