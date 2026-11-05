package rest

type getUpdateMovementFormView struct {
	ID            int64
	Title         string
	Merchant      string
	Amount        string
	Date          string
	MaxDate       string
	Methods       []getUpdateMovementFormOption
	Categories    []getUpdateMovementFormOptionGroup
	Note          string
	FromStatement bool
	Errors        []string
}

type getUpdateMovementFormOption struct {
	Value    string
	Label    string
	Color    string
	Selected bool
}

type getUpdateMovementFormOptionGroup struct {
	Label   string
	Options []getUpdateMovementFormOption
}
