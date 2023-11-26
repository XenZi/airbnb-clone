package domains

type RequestResetPassword struct {
	Email string `json:"email"`
	Token string `json:"token"`
}