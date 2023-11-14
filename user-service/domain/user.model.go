package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	Email     string             `json:"email" bson:"email"`
	Role      string             `json:"role" bson:"role"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Residence string             `json:"residence" bson:"residence"`
}

type UserDTO struct {
	ID        primitive.ObjectID `json:"id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Role      string             `json:"role"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Residence string             `json:"residence"`
}
