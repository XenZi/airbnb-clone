package domains

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SuccessfullyLoggedUser struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}
