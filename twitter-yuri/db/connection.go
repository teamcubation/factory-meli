package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() (*mongo.Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority&appName=%s", os.Getenv("MONGO_USERNAME"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_HOST"), os.Getenv("MONGO_APP_NAME"))
	fmt.Println("Connecting to MongoDB at:", uri)
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return client, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			return
		}
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return client, err

	}

	return client, nil
}

func ConnectToDatabase() (*mongo.Database, error) {
	client, err := Connect()
	if err != nil {
		return nil, err
	}

	database := client.Database("twitter")
	return database, nil
}

func ConnectToCollection(collectionName string) (*mongo.Collection, error) {
	database, err := ConnectToDatabase()
	if err != nil {
		return nil, err
	}

	collection := database.Collection(collectionName)
	return collection, nil
}
