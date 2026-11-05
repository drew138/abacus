package rest

type getUpdateMovementFormView struct {
	ID            int64                              `json:"id"`
	Title         string                             `json:"title"`
	Merchant      string                             `json:"merchant"`
	Amount        string                             `json:"amount"`
	Date          string                             `json:"date"`
	MaxDate       string                             `json:"max_date"`
	Methods       []getUpdateMovementFormOption      `json:"methods"`
	Categories    []getUpdateMovementFormOptionGroup `json:"categories"`
	Note          string                             `json:"note"`
	FromStatement bool                               `json:"from_statement"`
	Errors        []string                           `json:"errors"`
}

type getUpdateMovementFormOption struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	Color    string `json:"color"`
	Selected bool   `json:"selected"`
}

type getUpdateMovementFormOptionGroup struct {
	Label   string                        `json:"label"`
	Options []getUpdateMovementFormOption `json:"options"`
}
