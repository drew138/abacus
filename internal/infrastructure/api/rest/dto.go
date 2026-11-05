package rest

import (
	"time"

	"github.com/drew138/abacus/internal/domain"
)

type MovementForm struct {
	Merchant        string `form:"merchant" validate:"required,max=120"`
	Amount          int64  `form:"amount" validate:"required,gt=0,lte=1000000000"`
	Date            string `form:"date" validate:"required,datetime=2006-01-02"`
	Category        string `form:"category" validate:"required,max=60"`
	Method          string `form:"method" validate:"required,enum=method"`
	Note            string `form:"note" validate:"max=200"`
	ApplyToMerchant bool   `form:"apply_to_merchant"`
}

func (f MovementForm) toDomain(id int64) domain.Movement {
	date, _ := time.Parse("2006-01-02", f.Date)
	return domain.Movement{ID: id, Date: date, Merchant: f.Merchant, Amount: f.Amount, Category: f.Category, Method: domain.Method(f.Method), Note: f.Note}
}

type pageView struct {
	Title      string
	Sidebar    sidebarView
	Month      *monthView
	Trends     *trendsView
	Toast      *toastView
	CloseModal bool
}

type sidebarView struct {
	Trends bool
	Months []sidebarMonth
}

type sidebarMonth struct {
	Key      string
	Name     string
	Badge    int
	Amount   string
	Latest   bool
	Selected bool
}

type deltaView struct {
	Text  string
	Class string
}

type monthView struct {
	Key          string
	Title        string
	Spent        string
	Purchases    int
	HasPrevious  bool
	Delta        deltaView
	PreviousName string
	Donut        donutView
	Ranks        []rankView
	SavingsTotal string
	Savings      []savingView
	Expenses     expensesView
}

type donutView struct {
	Slices []sliceView
	Total  string
	Count  int
}

type sliceView struct {
	Key    string
	Name   string
	Color  string
	Path   string
	DX     string
	DY     string
	Short  string
	Share  string
	Filter string
	Amount string
}

type rankView struct {
	Name   string
	Color  string
	Amount string
	Width  string
	Share  string
	Delta  deltaView
	Sub    bool
	Filter string
	Hover  string
}

type savingView struct {
	Label  string
	Amount string
}

type expensesView struct {
	Month  string
	Count  int
	Filter string
	Search string
	Groups []expenseGroupView
	Empty  string
}

type expenseGroupView struct {
	Category  string
	Parent    string
	Color     string
	Count     int
	Total     string
	IsSavings bool
	Rows      []rowView
}

type rowView struct {
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

type trendsView struct {
	Empty     bool
	Subtitle  string
	Saved     string
	MonthsN   string
	Average   string
	LastName  string
	LastSpent string
	LastDelta deltaView
	Only      string
	Chart     barsView
	Legend    []legendView
	Months    []string
	Savings   tableRowView
	Rows      []tableRowView
	Totals    tableRowView
}

type legendView struct {
	Name  string
	Color string
	Off   bool
}

type tableRowView struct {
	Name  string
	Color string
	Cells []string
	Delta deltaView
	Subs  []tableRowView
}

type barsView struct {
	Width   int
	Height  int
	Grid    []gridView
	Columns []columnView
}

type gridView struct {
	Y     string
	Label string
}

type columnView struct {
	X        string
	LabelX   string
	TopY     string
	Total    string
	Month    string
	Segments []segmentView
}

type segmentView struct {
	Key   string
	Color string
	X     string
	Y     string
	W     string
	H     string
	Month string
	Value string
	Share string
	Delta deltaView
	Vs    string
}

type formView struct {
	ID            int64
	Title         string
	Merchant      string
	Amount        string
	Date          string
	MaxDate       string
	Method        string
	Methods       []optionView
	Categories    []optionGroupView
	Note          string
	FromStatement bool
	Errors        []string
}

type optionView struct {
	Value    string
	Label    string
	Color    string
	Selected bool
}

type optionGroupView struct {
	Label   string
	Options []optionView
}

type toastView struct {
	Message string
	UndoID  int64
}

type loginView struct {
	Error    string
	Username string
}

type ErrorResponse struct {
	Type   string        `json:"type"`
	Title  string        `json:"title"`
	Status int           `json:"status"`
	Errors []ErrorDetail `json:"errors,omitempty"`
}

type ErrorDetail struct {
	Code  string `json:"code"`
	Field string `json:"field,omitempty"`
	Param string `json:"param,omitempty"`
}
