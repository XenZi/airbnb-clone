package domains

type RegisterUser struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	CurrentPlace string `json:"currentPlace"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	Username     string `json:"username"`
	Age          int    `json:"age"`
}
