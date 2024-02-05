package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"user-service/config"
	"user-service/domain"
	"user-service/errors"
)

type UserRepository struct {
	cli    *mongo.Client
	logger *config.Logger
}

const (
	source  = "user-repository"
	userDB  = "user-service"
	userCol = "users"
)

func NewUserRepository(cli *mongo.Client, logger *config.Logger) *UserRepository {
	return &UserRepository{
		cli:    cli,
		logger: logger,
	}
}

func (ur UserRepository) CreatUser(user domain.User) (*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database(userDB).Collection(userCol)
	insertedUser, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		ur.logger.Println(err.Error())
		err, status := errors.HandleInsertError(err, user)
		if status == -1 {
			status = 500
		}
		ur.logger.LogError(source, err.Error())
		return nil, errors.NewError(err.Error(), status)
	}
	ur.logger.LogInfo(source, fmt.Sprintf("Inserted user by ID: %v", insertedUser))
	user.ID = insertedUser.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (ur UserRepository) GetAllUsers() ([]*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database(userDB).Collection(userCol)
	findOptions := options.Find()
	found, err := userCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		ur.logger.Println(err.Error())
		err, status := errors.HandleNoDocumentsError(err, "")
		if status == -1 {
			status = 500
		}
		ur.logger.LogError(source, err.Error())
		return nil, errors.NewError(err.Error(), status)
	}
	var users []*domain.User
	for found.Next(context.TODO()) {
		var user *domain.User
		err := found.Decode(&user)
		if err != nil {
			ur.logger.LogError(source, err.Error())
		}
		users = append(users, user)
	}
	ur.logger.LogInfo(source, fmt.Sprintf("Found %d users", len(users)))
	return users, nil
}

func (ur UserRepository) GetUserById(id string) (*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database(userDB).Collection(userCol)
	foundId, erro := primitive.ObjectIDFromHex(id)
	if erro != nil {
		ur.logger.LogError(source, erro.Error())
		return nil, errors.NewError(erro.Error(), 500)
	}
	filter := bson.D{{"_id", foundId}}
	var user *domain.User
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		err, status := errors.HandleNoDocumentsError(err, id)
		if status == -1 {
			status = 500
		}
		ur.logger.LogError(source, err.Error())
		return nil, errors.NewError(err.Error(), status)
	}
	ur.logger.LogInfo(source, fmt.Sprintf("Found user by id: %v ", user.ID.Hex()))
	return user, nil
}

func (ur UserRepository) UpdateUser(user domain.User) (*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database(userDB).Collection(userCol)
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{
		{"firstName", user.FirstName},
		{"lastName", user.LastName},
		{"residence", user.Residence},
		{"age", user.Age},
	}}}
	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		err, status := errors.HandleInsertError(err, user)
		ur.logger.LogError(source, err.Error())
		return nil, errors.NewError(err.Error(), status)
	}
	ur.logger.LogInfo(source, fmt.Sprintf("Updated user by id: %v ", user.ID.Hex()))
	return &user, nil
}

func (ur UserRepository) UpdateUserCreds(user domain.User) (*domain.User, *errors.ErrorStruct) {
	userCollection := ur.cli.Database(userDB).Collection(userCol)
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{
		{"username", user.Username},
		{"email", user.Email},
	}}}
	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		err, status := errors.HandleInsertError(err, user)
		ur.logger.LogError(source, err.Error())
		return nil, errors.NewError(err.Error(), status)
	}
	ur.logger.LogInfo(source, fmt.Sprintf("Updated user creds by id: %v ", user.ID.Hex()))
	return &user, nil
}

func (ur UserRepository) DeleteUser(id string) *errors.ErrorStruct {
	userCollection := ur.cli.Database(userDB).Collection(userCol)
	foundId, erro := primitive.ObjectIDFromHex(id)
	if erro != nil {
		ur.logger.LogError(source, erro.Error())
		return errors.NewError(erro.Error(), 500)
	}
	filter := bson.D{{"_id", foundId}}
	_, err := userCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		err, status := errors.HandleNoDocumentsError(err, id)
		ur.logger.LogError(source, err.Error())
		return errors.NewError(err.Error(), status)
	}
	ur.logger.LogInfo(source, fmt.Sprintf("Deleted user by id: %s ", id))
	return nil
}

func (ur UserRepository) UpdateRating(id string, rating float64) *errors.ErrorStruct {
	userCollection := ur.cli.Database(userDB).Collection(userCol)
	foundId, erro := primitive.ObjectIDFromHex(id)
	if erro != nil {
		return errors.NewError(erro.Error(), 500)
	}
	filter := bson.D{{"_id", foundId}}
	update := bson.D{{"$set", bson.D{
		{"rating", rating},
	}}}
	updatedId, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		err, status := errors.HandleNoDocumentsError(err, id)
		ur.logger.LogError(source, err.Error())
		return errors.NewError(err.Error(), status)
	}
	ur.logger.LogInfo(source, fmt.Sprintf("Updated user by id: %v ", updatedId))
	return nil
}
