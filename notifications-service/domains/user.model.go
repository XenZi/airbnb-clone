package domains

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName" `
	Residence string `json:"residence"`
	Age       int    `json:"age"`
	Rating    int    `json:"rating"`
}