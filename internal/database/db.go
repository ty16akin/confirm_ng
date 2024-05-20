package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	model "github.com/ty16akin/ConfirmNG/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
)

func Connect() {
	// Set client options
	godotenv.Load()
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set database and collection
	database = client.Database("confirm_ng")
	collection = database.Collection("users") // Change "users" to your collection name
}

func GetOne(u *model.User, filter interface{}) error {
	//Will automatically create a collection if not available
	collection := client.Database(database.Name()).Collection(collection.Name())
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(u)
	return err
}

func Get(filter interface{}) []*model.User {
	collection := client.Database(database.Name()).Collection(collection.Name())
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var result []*model.User
	for cur.Next(ctx) {
		user := &model.User{}
		er := cur.Decode(user)
		if er != nil {
			log.Fatal(er)
		}
		result = append(result, user)
	}
	return result
}

func AddOne(u *model.User) (*mongo.InsertOneResult, error) {
	collection := client.Database(database.Name()).Collection(collection.Name())
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, u)
	return result, err
}

func Update(u *model.User, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := client.Database(database.Name()).Collection(collection.Name())
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.UpdateMany(ctx, filter, update)
	return result, err
}

func RemoveOne(filter interface{}) (*mongo.DeleteResult, error) {
	collection := client.Database(database.Name()).Collection(collection.Name())
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}
