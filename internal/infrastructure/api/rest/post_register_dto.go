package rest

type RegisterForm struct {
	Username string `form:"username" validate:"required,max=60"`
	Password string `form:"password" validate:"required,min=12,max=72"`
	Confirm  string `form:"confirm" validate:"eqfield=Password"`
}

type postRegisterView struct {
	Error    string
	Username string
}
