package repository

import (
	do "accommodations-service/domain"
	"accommodations-service/errors"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type AccommodationRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

func NewAccommodationRepository(cli *mongo.Client, logger *log.Logger) *AccommodationRepo {
	return &AccommodationRepo{
		cli:    cli,
		logger: logger,
	}
}
func (ar *AccommodationRepo) SaveAccommodation(accommodation do.Accommodation) (*do.Accommodation, *errors.ErrorStruct) {
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	insertedAccommodation, err := accommodationCollection.InsertOne(context.TODO(), accommodation)
	if err != nil {
		ar.logger.Println(err.Error())

		return nil, errors.NewError(err.Error(), 500)
	}
	ar.logger.Println("Inserted ID is %v", insertedAccommodation)
	accommodation.Id = insertedAccommodation.InsertedID.(primitive.ObjectID)
	return &accommodation, nil
}

func (ar *AccommodationRepo) GetAccommodationById(id string) (*do.Accommodation, *errors.ErrorStruct) {
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	var accommodation *do.Accommodation
	accommId, _ := primitive.ObjectIDFromHex(id)
	err := accommodationCollection.FindOne(context.TODO(), bson.M{"_id": accommId}).Decode(&accommodation)
	if err != nil {
		return nil, errors.NewError(
			"Not able to retrieve data",
			500)
	}
	return accommodation, nil
}

func (ar *AccommodationRepo) GetAllAccommodations() ([]*do.Accommodation, *errors.ErrorStruct) {
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	var accommodations []*do.Accommodation

	cursor, err := accommodationCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, errors.NewError(
			"Not able to retrieve data",
			500)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, context.TODO())

	for cursor.Next(context.TODO()) {
		var accommodation do.Accommodation
		if err := cursor.Decode(&accommodation); err != nil {
			return nil, errors.NewError(
				"Error decoding data",
				500)
		}
		accommodations = append(accommodations, &accommodation)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.NewError(
			"Cursor error",
			500)
	}

	return accommodations, nil
}

func (ar *AccommodationRepo) UpdateAccommodationById(accommodation do.Accommodation) (*do.Accommodation, *errors.ErrorStruct) {
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")

	filter := bson.D{{Key: "_id", Value: accommodation.Id}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "location", Value: accommodation.Location},
			{Key: "name", Value: accommodation.Name},
			{Key: "conveniences", Value: accommodation.Conveniences},
			{Key: "minNumOfVisitors", Value: accommodation.MinNumOfVisitors},
			{Key: "maxNumOfVisitors", Value: accommodation.MaxNumOfVisitors},
		}},
	}

	_, err := accommodationCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ar.logger.Println(err)
		return nil, errors.NewError("Unable to update, database error", 500)
	}

	return &accommodation, nil
}

func (ar *AccommodationRepo) DeleteAccommodationById(id string) *errors.ErrorStruct {
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	accommId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": accommId}

	_, err := accommodationCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		ar.logger.Println(err)
		return errors.NewError("Unable to delete, database error", 500)
	}

	return nil
}
