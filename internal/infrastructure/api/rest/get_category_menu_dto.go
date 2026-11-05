package rest

type getCategoryMenuView struct {
	ID       int64                   `json:"id"`
	Merchant string                  `json:"merchant"`
	Month    string                  `json:"month"`
	Options  []getCategoryMenuOption `json:"options"`
}

type getCategoryMenuOption struct {
	Name    string `json:"name"`
	Color   string `json:"color"`
	Current bool   `json:"current"`
}
