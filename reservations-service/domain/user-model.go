package domain

type User struct {
	Id string `json: "string"`
}

func NewUser(id string) *User {
	return &User{
		Id: id,
	}
}
