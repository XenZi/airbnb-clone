package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"user-service/config"
	"user-service/domain"
	"user-service/errors"
)

type UserRepository struct {
	cli    *mongo.Client
	logger *config.Logger
}

func NewUserRepository(cli *mongo.Client, logger *config.Logger) *UserRepository {
	return &UserRepository{
		cli:    cli,
		logger: logger,
	}
}
func (ur UserRepository) CreatUser(user domain.User) (*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database("user-service").Collection("users")
	insertedUser, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		err, status := errors.HandleInsertError(err, user)
		if status == -1 {
			status = 500
		}
		ur.logger.Error("Error while inserting new user", log.Fields{
			"module": "database",
			"error":  err.Error(),
		})
		return nil, errors.NewError(err.Error(), status)
	}
	ur.logger.Infof("Successfully inserted user with email: " + user.Email)
	user.ID = insertedUser.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (ur UserRepository) GetAllUsers() ([]*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database("user-service").Collection("users")
	findOptions := options.Find()

	found, err := userCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		ur.logger.Println(err.Error())
		err, status := errors.HandleNoDocumentsError(err, "")
		if status == -1 {
			status = 500
		}
		return nil, errors.NewError(err.Error(), status)
	}
	var users []*domain.User
	for found.Next(context.TODO()) {
		var user *domain.User
		err := found.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur UserRepository) GetUserById(id string) (*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database("user-service").Collection("users")
	foundId, erro := primitive.ObjectIDFromHex(id)
	if erro != nil {
		return nil, errors.NewError(erro.Error(), 500)
	}
	filter := bson.D{{"_id", foundId}}
	var user *domain.User
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		ur.logger.Println(err.Error())
		err, status := errors.HandleNoDocumentsError(err, id)
		if status == -1 {
			status = 500
		}
		return nil, errors.NewError(err.Error(), status)
	}
	return user, nil
}

func (ur UserRepository) UpdateUser(user domain.User) (*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database("user-service").Collection("users")
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{
		{"firstName", user.FirstName},
		//{"username", user.Username},
		//{"email", user.Email},
		{"lastName", user.LastName},
		{"residence", user.Residence},
		{"age", user.Age},
	}}}
	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ur.logger.Println(err.Error())
		err, status := errors.HandleInsertError(err, user)
		return nil, errors.NewError(err.Error(), status)
	}
	return &user, nil
}

func (ur UserRepository) UpdateUserCreds(user domain.User) (*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database("user-service").Collection("users")
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{
		{"username", user.Username},
		{"email", user.Email},
	}}}
	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ur.logger.Println(err.Error())
		err, status := errors.HandleInsertError(err, user)
		return nil, errors.NewError(err.Error(), status)
	}
	return &user, nil
}

func (ur UserRepository) DeleteUser(id string) *errors.ErrorStruct {
	userCollection := ur.cli.Database("user-service").Collection("users")
	foundId, erro := primitive.ObjectIDFromHex(id)
	if erro != nil {
		return errors.NewError(erro.Error(), 500)
	}
	filter := bson.D{{"_id", foundId}}
	_, err := userCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return errors.NewError("User not found", 404)
	}
	return nil
}

func (ur UserRepository) UpdateRating(id string, rating float64) *errors.ErrorStruct {
	userCollection := ur.cli.Database("user-service").Collection("users")
	foundId, erro := primitive.ObjectIDFromHex(id)
	if erro != nil {
		return errors.NewError(erro.Error(), 500)
	}
	filter := bson.D{{"_id", foundId}}
	update := bson.D{{"$set", bson.D{
		{"rating", rating},
	}}}
	gu, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ur.logger.Println(err.Error())
		err, status := errors.HandleNoDocumentsError(err, id)
		return errors.NewError(err.Error(), status)
	}
	log.Println(gu)
	return nil
}
