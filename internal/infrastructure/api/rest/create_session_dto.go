package rest

type CreateSessionRequest struct {
	Username string `json:"username" form:"username" validate:"required,max=60"`
	Password string `json:"password" form:"password" validate:"required,max=72"`
}

type createSessionView struct {
	Error    string `json:"error"`
	Username string `json:"username"`
}
