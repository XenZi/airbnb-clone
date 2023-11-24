package domains

type RequestResetPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Password string `json:"password"`
	ConfirmedPassword string `json:"confirmedPassword"`
}