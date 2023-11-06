package domain

import (
	"github.com/google/uuid"
)

type User struct {
	Id uuid.UUID
}

func (u User) Equals(user User) bool {
	return u.Id.String() == user.Id.String()
}
