package rest

type RecategoriseMerchantRequest struct {
	Merchant string `json:"merchant" form:"merchant" validate:"required,max=120"`
	Category string `json:"category" form:"category" validate:"required,max=60"`
	Month    string `json:"month" form:"month" validate:"required,datetime=2006-01"`
}
