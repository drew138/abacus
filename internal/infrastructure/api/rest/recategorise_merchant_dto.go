package rest

type RecategoriseForm struct {
	Merchant string `form:"merchant" validate:"required,max=120"`
	Category string `form:"category" validate:"required,max=60"`
	Month    string `form:"month" validate:"required,datetime=2006-01"`
}
