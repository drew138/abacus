package rest

import (
	"time"

	"github.com/drew138/abacus/internal/domain"
)

type updateMovementFormView struct {
	ID            int64
	Title         string
	Merchant      string
	Amount        string
	Date          string
	MaxDate       string
	Methods       []updateMovementFormOption
	Categories    []updateMovementFormOptionGroup
	Note          string
	FromStatement bool
	Errors        []string
}

type updateMovementFormOption struct {
	Value    string
	Label    string
	Color    string
	Selected bool
}

type updateMovementFormOptionGroup struct {
	Label   string
	Options []updateMovementFormOption
}

type UpdateMovementRequest struct {
	Merchant        string `form:"merchant" validate:"required,max=120"`
	Amount          int64  `form:"amount" validate:"required,gt=0,lte=1000000000"`
	Date            string `form:"date" validate:"required,datetime=2006-01-02"`
	Category        string `form:"category" validate:"required,max=60"`
	Method          string `form:"method" validate:"required,enum=method"`
	Note            string `form:"note" validate:"max=200"`
	ApplyToMerchant bool   `form:"apply_to_merchant"`
}

func (f UpdateMovementRequest) toDomain(id int64) domain.Movement {
	date, _ := time.Parse("2006-01-02", f.Date)
	return domain.Movement{ID: id, Date: date, Merchant: f.Merchant, Amount: f.Amount, Category: f.Category, Method: domain.Method(f.Method), Note: f.Note}
}
