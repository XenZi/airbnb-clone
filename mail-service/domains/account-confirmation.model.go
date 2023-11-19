package domains

type AccountConfirmation struct {
	Email string `json: "email"`
	Token string `json: "token"`
}
