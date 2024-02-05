package repository

import (
	"accommodations-service/config"
	do "accommodations-service/domain"
	"accommodations-service/errors"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
)

type AccommodationRepo struct {
	cli    *mongo.Client
	logger *config.Logger
	tracer trace.Tracer
}

func NewAccommodationRepository(cli *mongo.Client, logger *config.Logger, tracer trace.Tracer) *AccommodationRepo {

	return &AccommodationRepo{
		cli:    cli,
		logger: logger,
		tracer: tracer,
	}
}
func (ar *AccommodationRepo) SaveAccommodation(ctx context.Context, accommodation do.Accommodation) (*do.Accommodation, *errors.ErrorStruct) {
	ctx, span := ar.tracer.Start(ctx, "AccommodationRepo.SaveAccommodation")
	defer span.End()
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	accommodation.Rating = 0.0

	insertedAccommodation, err := accommodationCollection.InsertOne(context.TODO(), accommodation)
	if err != nil {
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Failed to get accommodation by id in ApproveAccommodation func with id"))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))

		return nil, errors.NewError(err.Error(), 500)
	}
	ar.logger.Println("Inserted ID is %v", insertedAccommodation)
	ar.logger.LogInfo("accommodation-repo", fmt.Sprintf("Accommodation inserted successfully"))
	accommodation.Id = insertedAccommodation.InsertedID.(primitive.ObjectID)
	return &accommodation, nil
}

func (ar *AccommodationRepo) GetAccommodationById(ctx context.Context, id string) (*do.Accommodation, *errors.ErrorStruct) {
	ctx, span := ar.tracer.Start(ctx, "AccommodationRepo.GetAccommodationById")
	defer span.End()
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	var accommodation *do.Accommodation
	accommId, _ := primitive.ObjectIDFromHex(id)
	err := accommodationCollection.FindOne(context.TODO(), bson.M{"_id": accommId}).Decode(&accommodation)
	if err != nil {
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Failed to find one accommodation by id %s", accommId))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
		return nil, errors.NewError(
			"Not able to retrieve data",
			500)
	}
	ar.logger.LogInfo("accommodation-repo", fmt.Sprintf("Accommodation inserted successfully"))
	return accommodation, nil
}

func (ar *AccommodationRepo) FindAccommodationByIds(ctx context.Context, ids []string) ([]*do.Accommodation, *errors.ErrorStruct) {
	ctx, span := ar.tracer.Start(ctx, "AccommodationRepo.FindAccommodationByIds")
	defer span.End()
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	log.Println("Idevi za get su ", ids)

	// Convert string IDs to primitive.ObjectID
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ar.logger.LogError("accommodations-repo", fmt.Sprintf("Failed to turn id into primitive object  %s", id))
			ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
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
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Failed to find one accommodation by id "))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
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
			ar.logger.LogError("accommodations-repo", fmt.Sprintf("Unable to decode accommodations "))
			ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
			return nil, errors.NewError(
				"Not able to retrieve data",
				500)
		}
		accommodations = append(accommodations, &accommodation)
	}
	ar.logger.LogInfo("accommodation-repo", fmt.Sprintf("Successfully found accommodations"))
	log.Println("akomodacije su", accommodations)
	return accommodations, nil
}

func (ar *AccommodationRepo) GetAllAccommodations(ctx context.Context) ([]*do.Accommodation, *errors.ErrorStruct) {
	ctx, span := ar.tracer.Start(ctx, "AccommodationRepo.GetAllAccommodations")
	defer span.End()
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	var accommodations []*do.Accommodation

	cursor, err := accommodationCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Unable to find all accommodations"))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
		return nil, errors.NewError(
			"Not able to retrieve data",
			500)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			ar.logger.LogError("accommodations-repo", fmt.Sprintf("Error closing cursor "))
			ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
			log.Println("Error")
		}
	}(cursor, context.TODO())

	for cursor.Next(context.TODO()) {
		var accommodation do.Accommodation

		if err := cursor.Decode(&accommodation); err != nil {
			log.Println(accommodation)
			ar.logger.LogError("accommodations-repo", fmt.Sprintf("Error decoding accommodation  "))
			ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
			return nil, errors.NewError(
				"Error decoding data",
				500)
		}
		accommodations = append(accommodations, &accommodation)
	}

	if err := cursor.Err(); err != nil {
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Cursor error "))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
		return nil, errors.NewError(
			"Cursor error",
			500)
	}
	ar.logger.LogInfo("accommodation-repo", fmt.Sprintf("Successfully found accommodations"))
	return accommodations, nil
}

func (ar *AccommodationRepo) UpdateAccommodationById(ctx context.Context, accommodation do.Accommodation) (*do.Accommodation, *errors.ErrorStruct) {
	ctx, span := ar.tracer.Start(ctx, "AccommodationRepo.UpdateAccommodationById")
	defer span.End()
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
			{Key: "status", Value: accommodation.Status},
		}},
	}

	_, err := accommodationCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Unable to update accommodation with id %s", accommodation.Id))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
		ar.logger.Println(err)
		return nil, errors.NewError("Unable to update, database error", 500)
	}
	ar.logger.LogInfo("accommodation-repo", fmt.Sprintf("Successfully updated accommodation"))
	return &accommodation, nil
}

func (ar *AccommodationRepo) UpdateAccommodationStatus(accommodation do.Accommodation, accomId string) (*do.Accommodation, *errors.ErrorStruct) {
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	id, _ := primitive.ObjectIDFromHex(accomId)
	log.Println("STATUS U REPOU JE", accommodation)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: accommodation.Status},
		}},
	}

	_, err := accommodationCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Unable to update accommodation status of accommodation with id %s", accommodation.Id))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
		ar.logger.Println(err)
		return nil, errors.NewError("Unable to update, database error", 500)
	}
	ar.logger.LogInfo("accommodation-repo", fmt.Sprintf("Successfully updated accommodation status"))
	return &accommodation, nil
}
func (ar *AccommodationRepo) PutAccommodationRating(ctx context.Context, accommodationID string, rating float32) *errors.ErrorStruct {
	ctx, span := ar.tracer.Start(ctx, "AccommodationRepo.PutAccommodationRating")
	defer span.End()
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
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Unable to update accommodation rating of accommodation with id %s", accommodationID))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
		ar.logger.Println(err)
		return errors.NewError("Unable to update rating, database error", 500)
	}
	ar.logger.LogInfo("accommodation-repo", fmt.Sprintf("Successfully updated accommodation rating"))
	return nil
}

func (ar *AccommodationRepo) DeleteAccommodationById(ctx context.Context, id string) *errors.ErrorStruct {
	ctx, span := ar.tracer.Start(ctx, "AccommodationRepo.DeleteAccommodationById")
	defer span.End()
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	accommId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": accommId}

	_, err := accommodationCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Unable to delete accommodation with id %s", id))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
		ar.logger.Println(err)
		return errors.NewError("Unable to delete, database error", 500)
	}
	ar.logger.LogInfo("accommodation-repo", fmt.Sprintf("Successfully deleted accommodation with id %s", id))
	return nil
}

func (ar *AccommodationRepo) DeleteAccommodationsByUserId(ctx context.Context, id string) *errors.ErrorStruct {
	ctx, span := ar.tracer.Start(ctx, "AccommodationRepo.DeleteAccommodationsByUserId")
	defer span.End()
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	userId := id
	filter := bson.M{"userId": userId} // Assuming userId is the field representing the user ID

	result, err := accommodationCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Unable to delete multiple accommodations"))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
		ar.logger.Println(err)
		return errors.NewError("Unable to delete, database error", 500)
	}

	// Check the number of deleted documents if needed

	deletedCount := result.DeletedCount
	log.Println(deletedCount)
	ar.logger.LogInfo("accommodation-repo", fmt.Sprintf("Successfully deleted accommodation with id %s", id))
	return nil
}

func (ar *AccommodationRepo) SearchAccommodations(ctx context.Context, city, country string, numOfVisitors int, maxPrice int, conveniences []string) ([]do.Accommodation, *errors.ErrorStruct) {
	ctx, span := ar.tracer.Start(ctx, "AccommodationRepo.SearchAccommodations")
	defer span.End()
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	filter := bson.M{}

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
	ctx = context.TODO()

	// Apply the filter and retrieve accommodations
	cursor, err := accommodationCollection.Find(ctx, filter)
	if err != nil {
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Unable to find accommodations for searched components"))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
		return nil, errors.NewError("Unable to find accommodations, database error", 500)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			ar.logger.LogError("accommodations-repo", fmt.Sprintf("Cursor closed"))
			ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))

		}
	}(cursor, ctx)

	// Iterate through the results and decode them into accommodations slice
	for cursor.Next(ctx) {
		var accommodation do.Accommodation // Replace Accommodation with your struct type
		if err := cursor.Decode(&accommodation); err != nil {
			ar.logger.LogError("accommodations-repo", fmt.Sprintf("Unable to decode accommodation"))
			ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
			return nil, errors.NewError("Unable to decode accommodations,error", 500)
		}
		accommodations = append(accommodations, accommodation)
	}

	if err := cursor.Err(); err != nil {
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Cursor error"))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
		return nil, errors.NewError("Unable to find accommodations, database error", 500)
	}
	log.Println(accommodations)
	ar.logger.LogInfo("accommodation-repo", fmt.Sprintf("Successfully found accommodations filtered by search"))
	return accommodations, nil
}

func (ar *AccommodationRepo) PutAccommodationStatus(accommodationID string, status string) *errors.ErrorStruct {
	accommodationCollection := ar.cli.Database("accommodations-service").Collection("accommodations")
	id, _ := primitive.ObjectIDFromHex(accommodationID)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: status},
		}},
	}

	// Perform the update operation
	_, err := accommodationCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ar.logger.LogError("accommodations-repo", fmt.Sprintf("Unable to update accommodation status"))
		ar.logger.LogError("accommodation-repo", fmt.Sprintf("Error:"+err.Error()))
		ar.logger.Println(err)
		return errors.NewError("Unable to update rating, database error", 500)
	}
	ar.logger.LogInfo("accommodation-repo", fmt.Sprintf("Successfully updated accommodation status"))
	return nil
}
