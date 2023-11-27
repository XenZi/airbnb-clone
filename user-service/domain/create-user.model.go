package domain

type CreateUser struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Residence string `json:"residence"`
	Role      string `json:"role"`
	Username  string `json:"username"`
	Age       int    `json:"age"`
}
