package http

type UserRequest struct {
	Username string `json:"username" required:"true" minLength:"3" maxLength:"30"`
	Password string `json:"password" required:"true" minLength:"4" maxLength:"30"`
}
