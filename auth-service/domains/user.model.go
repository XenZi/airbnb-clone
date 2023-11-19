package domains

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `json: "username" bson:"username"`
	Password  string             `json: "password" bson:"password"`
	Email     string             `json: "email" bson:"email"`
	Role      string             `json: "role" bson:"role"`
	Confirmed bool               `json: "confirmed" bson:"confirmed"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
