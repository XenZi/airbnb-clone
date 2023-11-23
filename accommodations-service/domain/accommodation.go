package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accommodation struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId           string             `json:"userId" bson:"userId"`
	UserName         string             `json:"username" bson:"username"`
	Name             string             `json:"name" bson:"name"`
	Location         string             `json:"location" bson:"location"`
	Conveniences     string             `json:"conveniences" bson:"conveniences"`
	MinNumOfVisitors int                `json:"minNumOfVisitors" bson:"minNumOfVisitors"`
	MaxNumOfVisitors int                `json:"maxNumOfVisitors" bson:"maxNumOfVisitors"`
}

type AccommodationDTO struct {
	Id               string `json:"id"`
	UserId           string `json:"userId" `
	UserName         string `json:"username" `
	Name             string `json:"name" `
	Location         string `json:"location" `
	Conveniences     string `json:"conveniences" `
	MinNumOfVisitors int    `json:"minNumOfVisitors" `
	MaxNumOfVisitors int    `json:"maxNumOfVisitors" `
}
