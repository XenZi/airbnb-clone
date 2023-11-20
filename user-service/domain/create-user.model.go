package domain

type CreateUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Residence string `json:"residence"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	Username  string `json:"username"`
	Age       int    `json:"age"`
}
