package rest

import (
	"time"

	"github.com/drew138/abacus/internal/domain"
)

type createMovementFormView struct {
	ID            int64                           `json:"id"`
	Title         string                          `json:"title"`
	Merchant      string                          `json:"merchant"`
	Amount        string                          `json:"amount"`
	Date          string                          `json:"date"`
	MaxDate       string                          `json:"max_date"`
	Methods       []createMovementFormOption      `json:"methods"`
	Categories    []createMovementFormOptionGroup `json:"categories"`
	Note          string                          `json:"note"`
	FromStatement bool                            `json:"from_statement"`
	Errors        []string                        `json:"errors"`
}

type createMovementFormOption struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	Color    string `json:"color"`
	Selected bool   `json:"selected"`
}

type createMovementFormOptionGroup struct {
	Label   string                     `json:"label"`
	Options []createMovementFormOption `json:"options"`
}

type CreateMovementRequest struct {
	Merchant        string `json:"merchant" form:"merchant" validate:"required,max=120"`
	Amount          int64  `json:"amount" form:"amount" validate:"required,gt=0,lte=1000000000"`
	Date            string `json:"date" form:"date" validate:"required,datetime=2006-01-02"`
	Category        string `json:"category" form:"category" validate:"required,max=60"`
	Method          string `json:"method" form:"method" validate:"required,enum=method"`
	Note            string `json:"note" form:"note" validate:"max=200"`
	ApplyToMerchant bool   `json:"apply_to_merchant" form:"apply_to_merchant"`
}

func (f CreateMovementRequest) toDomain(id int64) domain.Movement {
	date, _ := time.Parse("2006-01-02", f.Date)
	return domain.Movement{ID: id, Date: date, Merchant: f.Merchant, Amount: f.Amount, Category: f.Category, Method: domain.Method(f.Method), Note: f.Note}
}
