package rest

type LoginForm struct {
	Username string `form:"username" validate:"required,max=60"`
	Password string `form:"password" validate:"required,max=72"`
}
