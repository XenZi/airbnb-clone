package repository

import (
	"log"

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