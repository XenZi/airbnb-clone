package repository

import (
	"auth-service/domains"
	"auth-service/errors"
	"context"
	errors2 "errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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

func (u UserRepository) SaveUser(user domains.User) (*domains.User, error) {
	userCollection := u.cli.Database("auth").Collection("user")
	insertedUser, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		u.logger.Println(err.Error())
		return nil, errors.HandleInsertError(err, user)
	}
	u.logger.Println("Inserted ID is %v", insertedUser)
	user.ID = insertedUser.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (u UserRepository) FindUserByEmail(email string) (*domains.User, error) {
	userCollection := u.cli.Database("auth").Collection("user")

	var user domains.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, errors2.New("Bad credentials")
	}
	return &user, nil
}
