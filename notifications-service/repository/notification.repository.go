package repository

import (
	"context"
	"notifications-service/config"
	"notifications-service/domains"
	"notifications-service/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
)

type NotificationRepository struct {
	cli    *mongo.Client
	logger *config.Logger
	tracer trace.Tracer
}

func NewNotificationRepository(cli *mongo.Client, logger *config.Logger, tracer trace.Tracer) *NotificationRepository {
	return &NotificationRepository{
		cli:    cli,
		logger: logger,
		tracer: tracer,
	}
}

func (nr NotificationRepository) CreateNewUserNotification(ctx context.Context, id string) (*domains.UserNotification, *errors.ErrorStruct) {
	ctx, span := nr.tracer.Start(ctx, "NotificationRepository.CreateNewUserNotification")
	defer span.End()
	userCollection := nr.cli.Database("notifications").Collection("notifications")
	primObjIdFromHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		nr.logger.LogError("notification-repository", err.Error())
		return nil, errors.NewError(err.Error(), 500)
	}
	notificationUser := domains.UserNotification{
		ID:            primObjIdFromHex,
		Notifications: []domains.Notification{},
	}
	_, err = userCollection.InsertOne(context.TODO(), notificationUser)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	nr.logger.LogInfo("notifications-repository", "Successfully created structure for new user notification.")
	return &notificationUser, nil
}

func (nr NotificationRepository) FindOneUserNotificationByID(ctx context.Context, id string) (*domains.UserNotification, *errors.ErrorStruct) {
	ctx, span := nr.tracer.Start(ctx, "NotificationRepository.FindOneUserNotificationByID")
	defer span.End()
	notificationCollection := nr.cli.Database("notifications").Collection("notifications")
	var userNotification domains.UserNotification
	primitiveObjectID, _ := primitive.ObjectIDFromHex(id)
	err := notificationCollection.FindOne(context.TODO(), bson.M{"_id": primitiveObjectID}).Decode(&userNotification)
	if err != nil {
		nr.logger.LogError("notification-repository", "Bad id "+id)
		return nil, errors.NewError(
			"Bad ID",
			401)
	}
	nr.logger.LogInfo("notifications-repository", "Found user notifications.")

	return &userNotification, nil

}

func (nr NotificationRepository) UpdateNotificationByID(ctx context.Context, userNotification *domains.UserNotification) (*domains.UserNotification, *errors.ErrorStruct) {
	ctx, span := nr.tracer.Start(ctx, "NotificationRepository.UpdateNotificationByID")
	defer span.End()
	userCollection := nr.cli.Database("notifications").Collection("notifications")
	filter := bson.D{{"_id", userNotification.ID}}

	if len(userNotification.Notifications) > 5 {
		userNotification.Notifications = userNotification.Notifications[:5]
		nr.logger.LogInfo("notifications-repository", "User had more than five notifications, deleting old.")
	}

	update := bson.M{
		"$set": bson.M{
			"notifications": userNotification.Notifications,
		},
	}

	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		nr.logger.LogError("notification-repository", err.Error())
		return nil, errors.NewError(err.Error(), 500)
	}

	nr.logger.LogInfo("notifications-repository", "Updating user notifications.")
	return userNotification, nil
}
