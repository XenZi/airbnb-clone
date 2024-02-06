package services

import (
	"context"
	"fmt"
	"notifications-service/config"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoService struct {
	cli    *mongo.Client
	logger *config.Logger
}

func New(ctx context.Context, logger *config.Logger) (*MongoService, error) {
	uri := os.Getenv("MONGO_DB_URI")
	logger.Println("URI: ", uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
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

	// Check connection -> if no error, connection is established
	err := s.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		s.logger.Println(err)
	}

	fmt.Println("Connection is valid")

}
func (s MongoService) GetCli() *mongo.Client {
	return s.cli
}
