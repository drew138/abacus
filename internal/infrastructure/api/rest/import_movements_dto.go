package rest

import (
	"time"

	"github.com/drew138/abacus/internal/domain"
)

type ImportRequest struct {
	Categories []ImportCategory `json:"categories" validate:"max=200,dive"`
	Movements  []ImportMovement `json:"movements" validate:"required,max=5000,dive"`
}

type ImportCategory struct {
	Name   string `json:"name" validate:"required,max=60"`
	Parent string `json:"parent" validate:"max=60"`
}

type ImportMovement struct {
	Ref      string `json:"ref" validate:"required,max=200"`
	Date     string `json:"date" validate:"required,datetime=2006-01-02"`
	Merchant string `json:"merchant" validate:"required,max=120"`
	Amount   int64  `json:"amount" validate:"required,gt=0,lte=1000000000"`
	Category string `json:"category" validate:"required,max=60"`
	Method   string `json:"method" validate:"required,enum=method"`
}

func (r ImportRequest) toDomain() ([]domain.Movement, []domain.Category) {
	categories := make([]domain.Category, 0, len(r.Categories))
	for _, c := range r.Categories {
		categories = append(categories, domain.Category{Name: c.Name, Parent: c.Parent})
	}
	movements := make([]domain.Movement, 0, len(r.Movements))
	for _, m := range r.Movements {
		date, _ := time.Parse("2006-01-02", m.Date)
		movements = append(movements, domain.Movement{Ref: m.Ref, Date: date, Merchant: m.Merchant, Amount: m.Amount, Category: m.Category, Method: domain.Method(m.Method)})
	}
	return movements, categories
}

type ImportResponse struct {
	Added   int `json:"added"`
	Updated int `json:"updated"`
	Kept    int `json:"kept"`
}
