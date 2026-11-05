package rest

type CreateUserRequest struct {
	Username string `json:"username" form:"username" validate:"required,max=60"`
	Password string `json:"password" form:"password" validate:"required,min=12,max=72"`
	Confirm  string `json:"confirm" form:"confirm" validate:"eqfield=Password"`
}

type createUserView struct {
	Error    string `json:"error"`
	Username string `json:"username"`
}
