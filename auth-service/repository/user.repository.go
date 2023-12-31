package repository

import (
	"auth-service/domains"
	"auth-service/errors"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	cli    *mongo.Client
	logger *log.Logger
}

func NewUserRepository(cli *mongo.Client, logger *log.Logger) *UserRepository {
	return &UserRepository{
		cli:    cli,
		logger: logger,
	}
}

func (u UserRepository) SaveUser(user domains.User) (*domains.User, *errors.ErrorStruct) {
	userCollection := u.cli.Database("auth").Collection("user")
	insertedUser, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		u.logger.Println(err.Error())
		err, status := errors.HandleInsertError(err, user)
		if status == -1 {
			status = 500
		}
		return nil, errors.NewError(err.Error(), status)
	}
	u.logger.Println("Inserted ID is %v", insertedUser)
	user.ID = insertedUser.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (u UserRepository) FindUserByEmail(email string) (*domains.User, *errors.ErrorStruct) {
	userCollection := u.cli.Database("auth").Collection("user")
	var user domains.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, errors.NewError(
			"Bad credentials",
			401)
	}
	log.Println(user)
	return &user, nil
}


func (u UserRepository) FindUserById(id string) (*domains.User, *errors.ErrorStruct) {
	userCollection := u.cli.Database("auth").Collection("user")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	var user domains.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, errors.NewError(
			"Not found with following ID",
			401)
	}
	return &user, nil
}

func (u UserRepository) UpdateUserConfirmation(id string) (*domains.User, *errors.ErrorStruct) {
	database := u.cli.Database("auth")
	collection := database.Collection("user")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewError(err.Error(),500)
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
		// Define the update to be applied
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "confirmed", Value: true},
			}},
		}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}

	if updateResult.ModifiedCount == 0 {
		return nil, errors.NewError("User not found or your account is already confirmed", 404)
	}

	user, errFromUserFinding := u.FindUserById(id)
	if err != nil {
		return nil, errFromUserFinding
	}
	return user, nil
}

func (u UserRepository) UpdateUserPassword(id string, newPassword string) (*domains.User, *errors.ErrorStruct) {
	database := u.cli.Database("auth")
	collection := database.Collection("user")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewError(err.Error(),500)
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
		return nil, errors.NewError(err.Error(), 500)
	}

	if updateResult.ModifiedCount == 0 {
		return nil, errors.NewError("User not found or your account", 400)
	}

	user, errFromUserFinding := u.FindUserById(id)
	if err != nil {
		return nil, errFromUserFinding
	}
	return user, nil
}

func (u UserRepository) UpdateUserCredentials(user domains.User) (*domains.User, *errors.ErrorStruct) {
	database := u.cli.Database("auth")
	collection := database.Collection("user")
	objectID, err := primitive.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		return nil, errors.NewError(err.Error(),500)
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
		return nil, errors.NewError(err.Error(), 500)
	}

	if updateResult.ModifiedCount == 0 {
		return nil, errors.NewError("User not found or your account", 400)
	}

	foundUser, errFromUserFinding := u.FindUserById(user.ID.Hex())
	if err != nil {
		return nil, errFromUserFinding
	}
	return foundUser, nil
}

func (u UserRepository) FindUserByUsername(username string) (*domains.User, *errors.ErrorStruct) {
	userCollection := u.cli.Database("auth").Collection("user")
	var user domains.User
	err := userCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, errors.NewError(
			"Not found with following ID",
			401)
	}
	return &user, nil
}


func (u UserRepository) DeleteUserById(id string) (*domains.User, *errors.ErrorStruct) {
	user, err := u.FindUserById(id)
	if err != nil {
		return nil, err
	}
	userCollection := u.cli.Database("auth").Collection("user")
	primitiveID, errFromCast := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewError(errFromCast.Error(), 500)
	}
	filter := bson.M{"_id": primitiveID}
	_, errFromDelete := userCollection.DeleteOne(context.TODO(), filter)
	if errFromDelete != nil {
		return nil, errors.NewError(errFromDelete.Error(), 500)
	}
	return user, nil
}