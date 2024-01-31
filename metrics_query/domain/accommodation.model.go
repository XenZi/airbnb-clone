package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Accommodation struct {
	Id                                 string            `json:"id"`
	ReportingDate                      time.Time         `json:"reportingDate" bson:"reportingDate"`
	OnScreenTime                       float64           `json:"onScreenTime" bson:"onScreenTime"`
	NumberOfVisits                     uint32            `json:"numberOfVisits" bson:"numberOfVisits"`
	NotClosedEventTimeStamps           map[string]string `json:"notClosedEventTimeStamps" bson:"notClosedEventTimeStamps"`
	LastAppliedUserJoinedEventNumber   int64             `json:"lastAppliedUserJoinedEventNumber" bson:"lastAppliedUserJoinedEventNumber"`
	LastAppliedUserLeftEventNumber     int64             `json:"lastAppliedUserLeftEventNumber" bson:"lastAppliedUserLeftEventNumber"`
	NumberOfReservations               uint32            `json:"numberOfReservations" bson:"numberOfReservations"`
	LastAppliedUserReservedEventNumber int64             `json:"lastAppliedUserReservedEventNumber" bson:"lastAppliedUserReservedEventNumber"`
	NumberOfRatings                    uint32            `json:"numberOfRatings" bson:"numberOfRatings"`
	LastAppliedUserRatedEventNumber    int64             `json:"lastAppliedUserRatedEventNumber" bson:"lastAppliedUserRatedEventNumber"`
}

type DBAccommodation struct {
	Id                                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ReportingDate                      time.Time          `json:"reportingDate" bson:"reportingDate"`
	OnScreenTime                       float64            `json:"onScreenTime" bson:"onScreenTime"`
	NumberOfVisits                     uint32             `json:"numberOfVisits" bson:"numberOfVisits"`
	NotClosedEventTimeStamps           map[string]string  `json:"notClosedEventTimeStamps" bson:"notClosedEventTimeStamps"`
	LastAppliedUserJoinedEventNumber   int64              `json:"lastAppliedUserJoinedEventNumber" bson:"lastAppliedUserJoinedEventNumber"`
	LastAppliedUserLeftEventNumber     int64              `json:"lastAppliedUserLeftEventNumber" bson:"lastAppliedUserLeftEventNumber"`
	NumberOfReservations               uint32             `json:"numberOfReservations" bson:"numberOfReservations"`
	LastAppliedUserReservedEventNumber int64              `json:"lastAppliedUserReservedEventNumber" bson:"lastAppliedUserReservedEventNumber"`
	NumberOfRatings                    uint32             `json:"numberOfRatings" bson:"numberOfRatings"`
	LastAppliedUserRatedEventNumber    int64              `json:"lastAppliedUserRatedEventNumber" bson:"lastAppliedUserRatedEventNumber"`
}
