package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() (*mongo.Client, context.Context, context.CancelFunc, error) {
	MONGODB_URL := os.Getenv("MONGO_DB_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, connectErr := mongo.Connect(ctx, options.Client().ApplyURI(MONGODB_URL))
	if connectErr != nil {
		cancel()
		log.Fatal("Error connecting DB", connectErr)
		return nil, nil, nil, connectErr

	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
		client.Disconnect(ctx)
		return nil, nil, nil, err
	}

	DB = client.Database("bug-point")
	fmt.Println("Connected to DB")
	return client, ctx, cancel, connectErr
}

func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}
