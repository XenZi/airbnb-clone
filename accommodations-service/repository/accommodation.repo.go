package repository

import (
	"context"
	do "github.com/XenZi/airbnb-clone/accommodations-service/domain"
	"github.com/XenZi/airbnb-clone/accommodations-service/errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	accommodation.Rating = 0.0

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

func (ar *AccommodationRepo) FindAccommodationByIds(ids []string) ([]*do.Accommodation, *errors.ErrorStruct) {
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	log.Println("Idevi za get su ", ids)

	// Convert string IDs to primitive.ObjectID
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.NewError(
				"Not able to convert to primitive",
				500) // Handle invalid ID error
		}
		objectIDs = append(objectIDs, objectID)
	}

	// Prepare the filter for finding accommodations by IDs
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}

	// Find accommodations
	cursor, err := accommodationCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, errors.NewError(
			"Not able to find data",
			500)
	}
	defer cursor.Close(context.TODO())

	// Decode results into Accommodation
	var accommodations []*do.Accommodation
	for cursor.Next(context.TODO()) {
		var accommodation do.Accommodation
		err := cursor.Decode(&accommodation)
		if err != nil {
			return nil, errors.NewError(
				"Not able to retrieve data",
				500)
		}
		accommodations = append(accommodations, &accommodation)
	}
	log.Println("akomodacije su", accommodations)
	return accommodations, nil
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
			log.Println(accommodation)
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
			{Key: "address", Value: accommodation.Address},
			{Key: "city", Value: accommodation.City},

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
func (ar *AccommodationRepo) PutAccommodationRating(accommodationID string, rating float32) *errors.ErrorStruct {
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	id, _ := primitive.ObjectIDFromHex(accommodationID)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "rating", Value: rating},
		}},
	}

	// Perform the update operation
	_, err := accommodationCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ar.logger.Println(err)
		return errors.NewError("Unable to update rating, database error", 500)
	}

	return nil
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

func (ar *AccommodationRepo) DeleteAccommodationsByUserId(id string) *errors.ErrorStruct {
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	userId := id
	filter := bson.M{"userId": userId} // Assuming userId is the field representing the user ID

	result, err := accommodationCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		ar.logger.Println(err)
		return errors.NewError("Unable to delete, database error", 500)
	}

	// Check the number of deleted documents if needed

	deletedCount := result.DeletedCount
	log.Println(deletedCount)

	return nil
}

func (ar *AccommodationRepo) SearchAccommodations(city, country string, numOfVisitors int, minPrice int, maxPrice int, conveniences []string) ([]do.Accommodation, *errors.ErrorStruct) {
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	filter := bson.M{}
	log.Println("minimalnaCijena", minPrice)
	log.Println("maximalnaCijena", maxPrice)
	log.Println("convies", conveniences)
	// Build the filter based on the provided parameters
	if city != "" {
		filter["city"] = city
	}
	if country != "" {
		filter["country"] = country
	}

	if numOfVisitors > 0 {
		filter["$and"] = bson.A{
			bson.M{"minNumOfVisitors": bson.M{"$lte": numOfVisitors}},
			bson.M{"maxNumOfVisitors": bson.M{"$gte": numOfVisitors}},
		}
	}

	if numOfVisitors > 0 {
		filter["$and"] = bson.A{
			bson.M{"minNumOfVisitors": bson.M{"$lte": numOfVisitors}},
			bson.M{"maxNumOfVisitors": bson.M{"$gte": numOfVisitors}},
		}
	}
	log.Println("duzina je", len(conveniences))
	if len(conveniences) > 0 {
		filter["conveniences"] = bson.M{"$in": conveniences}
	}

	// Perform the search using the constructed filter
	var accommodations []do.Accommodation // Replace Accommodation with your struct type
	ctx := context.TODO()

	// Apply the filter and retrieve accommodations
	cursor, err := accommodationCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.NewError("Unable to find accommodations, database error", 500)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	// Iterate through the results and decode them into accommodations slice
	for cursor.Next(ctx) {
		var accommodation do.Accommodation // Replace Accommodation with your struct type
		if err := cursor.Decode(&accommodation); err != nil {
			return nil, errors.NewError("Unable to decode accommodations,error", 500)
		}
		accommodations = append(accommodations, accommodation)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.NewError("Unable to find accommodations, database error", 500)
	}
	log.Println(accommodations)

	return accommodations, nil
}
