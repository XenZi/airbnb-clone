package services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

type MongoService struct {
	cli    *mongo.Client
	logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger) (*MongoService, error) {
	uri := os.Getenv("MONGO_DB_URI")
	logger.Println("URI: ", uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	indexName, err := client.Database("auth").Collection("user").Indexes().CreateMany(context.Background(),
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "username", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		})
	logger.Println("Constraints created for mongo-db Auth model with name ", indexName)
	if err != nil {
		logger.Fatal("Error while connecting to auth-mongo", err)
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
