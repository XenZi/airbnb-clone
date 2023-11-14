package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"user-service/domain"
	"user-service/errors"
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
func (ur UserRepository) CreatUser(user domain.User) (*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database("users").Collection("user")
	insertedUser, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		ur.logger.Println(err.Error())
		err, status := errors.HandleInsertError(err, user)
		if status == -1 {
			status = 500
		}
		return nil, errors.NewError(err.Error(), status)
	}
	ur.logger.Println("Inserted ID is %v", insertedUser)
	user.ID = insertedUser.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (ur UserRepository) GetAllUsers() ([]*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database("users").Collection("user")
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
	userCollection := ur.cli.Database("users").Collection("user")
	filter := bson.D{{"id", id}}
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
	userCollection := ur.cli.Database("users").Collection("user")
	filter := bson.D{{"id", user.ID}}
	update := bson.D{{"$set", bson.D{{"username", user.Username}}}}
	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ur.logger.Println(err.Error())
		err, status := errors.HandleInsertError(err, user)
		if status == -1 {
			status = 500
		}
		return nil, errors.NewError(err.Error(), status)
	}
	return &user, nil
}
