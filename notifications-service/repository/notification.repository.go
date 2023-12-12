package repository

import (
	"context"
	"log"
	"notifications-service/domains"
	"notifications-service/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository struct {
	cli    *mongo.Client
	logger *log.Logger
}

func NewNotificationRepository(cli *mongo.Client, logger *log.Logger) *NotificationRepository {
	return &NotificationRepository{
		cli: cli,
		logger: logger,
	}
}

func (nr NotificationRepository) CreateNewUserNotification(id string) (*domains.UserNotification, *errors.ErrorStruct){
	userCollection := nr.cli.Database("notifications").Collection("notifications")
	primObjIdFromHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	notificationUser := domains.UserNotification{
		ID: primObjIdFromHex,
		Notifications: []domains.Notification{},
	}
	_, err = userCollection.InsertOne(context.TODO(), notificationUser)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	return &notificationUser, nil
}

func (nr NotificationRepository) FindOneUserNotificationByID(id string) (*domains.UserNotification, *errors.ErrorStruct) {
	notificationCollection := nr.cli.Database("notifications").Collection("notifications")
	var userNotification domains.UserNotification
	primitiveObjectID, _ := primitive.ObjectIDFromHex(id)
	err := notificationCollection.FindOne(context.TODO(), bson.M{"_id": primitiveObjectID}).Decode(&userNotification)
	if err != nil {
		return nil, errors.NewError(
			"Bad ID",
			401)
	}
	return &userNotification, nil

}

func (nr NotificationRepository) UpdateNotificationByID(userNotification *domains.UserNotification) (*domains.UserNotification, *errors.ErrorStruct) {
	userCollection := nr.cli.Database("notifications").Collection("notifications")
	filter := bson.D{{"_id", userNotification.ID}}
	update := bson.D{{"$set", bson.D{
		{"notifications", userNotification.Notifications},
	}}}
	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	return userNotification, nil
}