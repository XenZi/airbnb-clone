package repository

import (
	"auth-service/config"
	"auth-service/domains"
	"auth-service/errors"
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
)

type UserRepository struct {
	cli    *mongo.Client
	logger *config.Logger
	tracer trace.Tracer
}

func NewUserRepository(cli *mongo.Client, logger *config.Logger, tracer trace.Tracer) *UserRepository {
	return &UserRepository{
		cli:    cli,
		logger: logger,
		tracer: tracer,
	}
}

func (u UserRepository) SaveUser(ctx context.Context, user domains.User) (*domains.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.SaveUser")
	defer span.End()
	userCollection := u.cli.Database("auth").Collection("user")
	insertedUser, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		err, status := errors.HandleInsertError(err, user)
		if status == -1 {
			status = 500
		}
		u.logger.Error("Error while inserting new user", log.Fields{
			"module": "database",
			"error":  err.Error(),
		})
		return nil, errors.NewError(err.Error(), status)
	}
	u.logger.Infof("Successfully inserted user with email " + user.Email)
	user.ID = insertedUser.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (u UserRepository) FindUserByEmail(ctx context.Context, email string) (*domains.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserRepository.FindUserByEmail")
	defer span.End()
	userCollection := u.cli.Database("auth").Collection("user")
	var user domains.User
	u.logger.Infof("Looking for a user with email " + email)
	err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errors.NewError(
			"Bad credentials",
			401)
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User found with email %v", email))
	return &user, nil
}

func (u UserRepository) FindUserById(ctx context.Context, id string) (*domains.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.FindUserByID")
	defer span.End()
	userCollection := u.cli.Database("auth").Collection("user")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errors.NewError(err.Error(), 500)
	}
	var user domains.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		u.logger.LogError("auth-db", "Not found with following ID")
		return nil, errors.NewError(
			"Not found with following ID",
			401)
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User found by id %v", id))
	return &user, nil
}

func (u UserRepository) UpdateUserConfirmation(ctx context.Context, id string) (*domains.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.UpdateUserConfirmation")
	defer span.End()
	database := u.cli.Database("auth")
	collection := database.Collection("user")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errors.NewError(err.Error(), 500)
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "confirmed", Value: true},
		}},
	}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errors.NewError(err.Error(), 500)
	}

	if updateResult.ModifiedCount == 0 {
		u.logger.LogError("auth-db", "User not found or your account is already confirmed")
		return nil, errors.NewError("User not found or your account is already confirmed", 404)
	}

	user, errFromUserFinding := u.FindUserById(ctx, id)
	u.logger.LogInfo("user-service", fmt.Sprintf("Updated user ID %v for confirmation", id))

	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errFromUserFinding
	}
	return user, nil
}

func (u UserRepository) UpdateUserPassword(ctx context.Context, id string, newPassword string) (*domains.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.UpdateUserPassword")
	defer span.End()
	database := u.cli.Database("auth")
	collection := database.Collection("user")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errors.NewError(err.Error(), 500)
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	// Define the update to be applied
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "password", Value: newPassword},
		}},
	}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errors.NewError(err.Error(), 500)
	}

	if updateResult.ModifiedCount == 0 {
		u.logger.LogError("auth-db", err.Error())
		return nil, errors.NewError("User not found or your account", 400)
	}

	user, errFromUserFinding := u.FindUserById(ctx, id)
	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errFromUserFinding
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("Updated passwor for user ID %v", id))

	return user, nil
}

func (u UserRepository) UpdateUserCredentials(ctx context.Context, user domains.User) (*domains.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.UpdateUserCredentials")
	defer span.End()
	database := u.cli.Database("auth")
	collection := database.Collection("user")
	objectID, err := primitive.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errors.NewError(err.Error(), 500)
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	// Define the update to be applied
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "email", Value: user.Email},
			{Key: "username", Value: user.Username},
		}},
	}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errors.NewError(err.Error(), 500)
	}

	if updateResult.ModifiedCount == 0 {
		u.logger.LogError("auth-db", err.Error())
		return nil, errors.NewError("User not found or your account", 400)
	}

	foundUser, errFromUserFinding := u.FindUserById(ctx, user.ID.Hex())
	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errFromUserFinding
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("Updated credentials for user ID %v", user.ID.Hex()))
	return foundUser, nil
}

func (u UserRepository) FindUserByUsername(ctx context.Context, username string) (*domains.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.FindUserByUsername")
	defer span.End()
	userCollection := u.cli.Database("auth").Collection("user")
	var user domains.User
	err := userCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		u.logger.LogError("auth-db", err.Error())
		return nil, errors.NewError(
			"Not found with following ID",
			401)
	}
	return &user, nil
}

func (u UserRepository) DeleteUserById(ctx context.Context, id string) (*domains.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.FindUserByUsername")
	defer span.End()
	user, err := u.FindUserById(ctx, id)
	if err != nil {
		u.logger.LogError("auth-db", err.GetErrorMessage())
		return nil, err
	}
	userCollection := u.cli.Database("auth").Collection("user")
	primitiveID, errFromCast := primitive.ObjectIDFromHex(id)
	if err != nil {
		u.logger.LogError("auth-db", errFromCast.Error())
		return nil, errors.NewError(errFromCast.Error(), 500)
	}
	filter := bson.M{"_id": primitiveID}
	_, errFromDelete := userCollection.DeleteOne(context.TODO(), filter)
	if errFromDelete != nil {
		u.logger.LogError("auth-db", errFromDelete.Error())
		return nil, errors.NewError(errFromDelete.Error(), 500)
	}
	return user, nil
}
