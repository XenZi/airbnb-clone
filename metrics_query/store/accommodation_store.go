package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"metrics_query/domain"
)

type AccommodationStore struct {
	cli *mongo.Client
}

func NewStore(cli *mongo.Client) domain.AccommodationStore {
	return AccommodationStore{cli: cli}
}
func (a AccommodationStore) Create(accommodation domain.Accommodation, collection string) error {
	db := a.cli.Database("accommodation-service").Collection(collection)
	acc, err := accToDBAcc(accommodation)
	if err != nil {
		log.Println("OVDE JE PRVI")
		return err
	}
	_, err = db.InsertOne(context.TODO(), acc)
	return nil
}

func (a AccommodationStore) Read(id, collection string) (*domain.Accommodation, error) {
	db := a.cli.Database("accommodation-service").Collection(collection)
	var accommodation *domain.DBAccommodation
	dbId, _ := primitive.ObjectIDFromHex(id)
	err := db.FindOne(context.TODO(), bson.M{"_id": dbId}).Decode(&accommodation)
	if err != nil {
		return nil, err
	}
	return dbToAcc(*accommodation), nil
}

func (a AccommodationStore) Update(accommodation domain.Accommodation, collection string) error {
	db := a.cli.Database("accommodation-service").Collection(collection)
	acc, err := accToDBAcc(accommodation)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: acc.Id}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{"reportingDate", acc.ReportingDate},
			{"onScreenTime", acc.OnScreenTime},
			{"numberOfVisits", acc.NumberOfVisits},
			{"notClosedEventTimeStamps", acc.NotClosedEventTimeStamps},
			{"lastAppliedUserJoinedEventNumber", acc.LastAppliedUserJoinedEventNumber},
			{"lastAppliedUserLeftEventNumber", acc.LastAppliedUserLeftEventNumber},
			{"numberOfReservations", acc.NumberOfReservations},
			{"lastAppliedUserReservedEventNumber", acc.LastAppliedUserReservedEventNumber},
			{"numberOfRatings", acc.NumberOfRatings},
			{"lastAppliedUserRatedEventNumber", acc.LastAppliedUserRatedEventNumber},
		}},
	}
	_, err = db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func accToDBAcc(acc domain.Accommodation) (*domain.DBAccommodation, error) {

	id, err := primitive.ObjectIDFromHex(acc.Id)
	if err != nil {
		log.Println("OVDE JE DRUGI")
		return nil, err
	}
	ret := domain.DBAccommodation{
		Id:                                 id,
		ReportingDate:                      acc.ReportingDate,
		OnScreenTime:                       acc.OnScreenTime,
		NumberOfVisits:                     acc.NumberOfVisits,
		NotClosedEventTimeStamps:           acc.NotClosedEventTimeStamps,
		LastAppliedUserJoinedEventNumber:   acc.LastAppliedUserJoinedEventNumber,
		LastAppliedUserLeftEventNumber:     acc.LastAppliedUserLeftEventNumber,
		NumberOfReservations:               acc.NumberOfReservations,
		LastAppliedUserReservedEventNumber: acc.LastAppliedUserReservedEventNumber,
		NumberOfRatings:                    acc.NumberOfRatings,
		LastAppliedUserRatedEventNumber:    acc.LastAppliedUserRatedEventNumber,
	}
	return &ret, nil
}

func dbToAcc(acc domain.DBAccommodation) *domain.Accommodation {
	ret := domain.Accommodation{
		Id:                                 acc.Id.Hex(),
		ReportingDate:                      acc.ReportingDate,
		OnScreenTime:                       acc.OnScreenTime,
		NumberOfVisits:                     acc.NumberOfVisits,
		NotClosedEventTimeStamps:           acc.NotClosedEventTimeStamps,
		LastAppliedUserJoinedEventNumber:   acc.LastAppliedUserJoinedEventNumber,
		LastAppliedUserLeftEventNumber:     acc.LastAppliedUserLeftEventNumber,
		NumberOfReservations:               acc.NumberOfReservations,
		LastAppliedUserReservedEventNumber: acc.LastAppliedUserReservedEventNumber,
		NumberOfRatings:                    acc.NumberOfRatings,
		LastAppliedUserRatedEventNumber:    acc.LastAppliedUserRatedEventNumber,
	}
	return &ret
}
