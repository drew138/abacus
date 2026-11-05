package rest

type categoryMenuView struct {
	ID       int64
	Merchant string
	Month    string
	Options  []categoryMenuOption
}

type categoryMenuOption struct {
	Name    string
	Color   string
	Current bool
}
