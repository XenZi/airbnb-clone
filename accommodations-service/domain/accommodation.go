package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accommodation struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId           string             `json:"userId" bson:"userId"`
	UserName         string             `json:"username" bson:"username"`
	Name             string             `json:"name" bson:"name"`
	Address          string             `json:"address" bson:"address"`
	City             string             `json:"city" bson:"city"`
	Country          string             `json:"country" bson:"country"`
	Conveniences     []string           `json:"conveniences" bson:"conveniences"`
	MinNumOfVisitors int                `json:"minNumOfVisitors" bson:"minNumOfVisitors"`
	MaxNumOfVisitors int                `json:"maxNumOfVisitors" bson:"maxNumOfVisitors"`
}

type CreateAccommodation struct {
	Id                          primitive.ObjectID            `bson:"_id,omitempty" json:"id"`
	UserId                      string                        `json:"userId" bson:"userId"`
	UserName                    string                        `json:"username" bson:"username"`
	Name                        string                        `json:"name" bson:"name"`
	Address                     string                        `json:"address" bson:"address"`
	City                        string                        `json:"city" bson:"city"`
	Country                     string                        `json:"country" bson:"country"`
	Conveniences                []string                      `json:"conveniences" bson:"conveniences"`
	MinNumOfVisitors            int                           `json:"minNumOfVisitors" bson:"minNumOfVisitors"`
	MaxNumOfVisitors            int                           `json:"maxNumOfVisitors" bson:"maxNumOfVisitors"`
	AvailableAccommodationDates []AvailableAccommodationDates `json:"availableAccommodationDates"`
	Location                    string                        `json:"location" `
}

type AvailableAccommodationDates struct {
	AccommodationId string `json:"id"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	Location        string `json:"location"`
	Price           int    `json:"price"`
}

type AccommodationDTO struct {
	Id               string   `json:"id"`
	UserId           string   `json:"userId" `
	UserName         string   `json:"username" `
	Name             string   `json:"name" `
	Address          string   `json:"address" `
	City             string   `json:"city" `
	Country          string   `json:"country" `
	Conveniences     []string `json:"conveniences" `
	MinNumOfVisitors int      `json:"minNumOfVisitors" `
	MaxNumOfVisitors int      `json:"maxNumOfVisitors" `
}
