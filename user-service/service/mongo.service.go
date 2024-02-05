package service

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
	"user-service/config"
)

type MongoService struct {
	cli    *mongo.Client
	logger *config.Logger
}

const mongoSource = "mongo-service"

func New(ctx context.Context, logger *config.Logger) (*MongoService, error) {
	uri := os.Getenv("MONGO_DB_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.LogError(mongoSource, err.Error())
		return nil, err
	}
	return &MongoService{
		cli:    client,
		logger: logger,
	}, nil
}

func (s MongoService) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		s.logger.LogError(mongoSource, err.Error())
	}
}
func (s MongoService) GetCli() *mongo.Client {
	return s.cli
}
