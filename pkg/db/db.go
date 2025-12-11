package db

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DbConnect(uri string)  *mongo.Client {

	client,err := mongo.Connect(options.Client().ApplyURI(uri))
	slog.Info("Connecting to MongoDB...")
	 
	if err != nil {
		panic(err)
	}
	return client
}

func GetCollection(client *mongo.Client, dbName string, collName string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collName)
	return collection
}

