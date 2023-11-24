package domains

type ChangePassword struct {
	OldPassword       string `json:"oldPassword"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmedPassword"`
}