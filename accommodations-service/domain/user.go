package domain

type User struct {
	Id       string
	username string
}

func (u User) Equals(user User) bool {
	return u.Id == user.Id
}
