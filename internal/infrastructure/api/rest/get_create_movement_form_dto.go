package rest

type getCreateMovementFormView struct {
	ID            int64
	Title         string
	Merchant      string
	Amount        string
	Date          string
	MaxDate       string
	Methods       []getCreateMovementFormOption
	Categories    []getCreateMovementFormOptionGroup
	Note          string
	FromStatement bool
	Errors        []string
}

type getCreateMovementFormOption struct {
	Value    string
	Label    string
	Color    string
	Selected bool
}

type getCreateMovementFormOptionGroup struct {
	Label   string
	Options []getCreateMovementFormOption
}
