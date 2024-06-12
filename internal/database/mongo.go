package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	Users  *mongo.Collection
)

func Init(uri string, database string) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI).SetServerSelectionTimeout(30 * time.Second)

	localClient, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return err
	}

	client = localClient

	Users = client.Database(database).Collection("users")

	err = client.Database("confirm_ng").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err()
	return err
}

func Close() error {
	return client.Disconnect(context.Background())
}
