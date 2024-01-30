package create_accommodation

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateAccommodation struct {
	Id                          primitive.ObjectID            `bson:"_id,omitempty" json:"id"`
	UserId                      string                        `json:"userId" bson:"userId"`
	UserName                    string                        `json:"username" bson:"username"`
	Email                       string                        `json:"email" bson:"email"`
	Name                        string                        `json:"name" bson:"name"`
	Address                     string                        `json:"address" bson:"address"`
	City                        string                        `json:"city" bson:"city"`
	Country                     string                        `json:"country" bson:"country"`
	Conveniences                []string                      `json:"conveniences" bson:"conveniences"`
	MinNumOfVisitors            int                           `json:"minNumOfVisitors" bson:"minNumOfVisitors"`
	MaxNumOfVisitors            int                           `json:"maxNumOfVisitors" bson:"maxNumOfVisitors"`
	AvailableAccommodationDates []AvailableAccommodationDates `json:"availableAccommodationDates"`
	Location                    string                        `json:"location" `
	Status                      AccommodationStatus           `json:"status"`
}

type AvailableAccommodationDates struct {
	AccommodationId string   `json:"accommodationId"`
	DateRange       []string `json:"dateRange"`
	Location        string   `json:"location"`
	Price           int      `json:"price"`
}

type AccommodationStatus string

const (
	Pending AccommodationStatus = "Pending"
	Created AccommodationStatus = "Created"
)

type CreateAccommodationCommandType int8

const (
	UpdateAccommodation CreateAccommodationCommandType = iota
	RollbackAccommodation
	ApproveAccommodation
	DeleteAccommodation
	SaveAvailabilities
	UnknownCommand
)

type CreateAccommodationCommand struct {
	Accommodation CreateAccommodation
	Type          CreateAccommodationCommandType
}

type CreateAccommodationReplyType int8

const (
	AccommodationCreated CreateAccommodationReplyType = iota
	AccommodationNotCreated
	AccommodationRolledBack
	AvailabilitiesCreated
	AvailabilitiesNotCreated
	AccommodationApproved
	AccommodationCancelled
	UnknownReply
)

type CreateAccommodationReply struct {
	Accommodation CreateAccommodation
	Type          CreateAccommodationReplyType
}
