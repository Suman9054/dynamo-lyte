package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Query struct {
	Querykey   string
	Queryvalue string
}

type Document struct {
	tanentId   string
	collection string
	partitionkey  string
	data       interface{}
}

func DbConnect(uri string) ( *mongo.Client,error) {

	client,err := mongo.Connect(options.Client().ApplyURI(uri))
	return client, err
}

func GetCollection(client *mongo.Client, dbName string, collName string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collName)
	return collection
}

func Insertusercollectionvalus(collection *mongo.Collection, document Document) (*mongo.InsertOneResult, error) {
	result,err:= collection.InsertOne(context.TODO(),bson.D{
		{Key: "tanentId", Value: document.tanentId},
		{Key: "collection", Value: document.collection},
		{Key: "partitionkey", Value: document.partitionkey},
		{Key: "data", Value: document.data},
	})
	return result, err
}


func FindDocument(collection *mongo.Collection, filter *Query) (*mongo.SingleResult) {
	result := collection.FindOne(context.TODO(), bson.D{{Key: filter.Querykey, Value: filter.Queryvalue}})
	return result
}

func UpdateDocument(collection *mongo.Collection, filter *Query, update interface{}) (*mongo.UpdateResult, error) {
	result, err := collection.UpdateOne(context.TODO(), bson.D{{Key: filter.Querykey, Value: filter.Queryvalue}}, update)
	return result, err
}

func DeleteDocument(collection *mongo.Collection, filter *Query) (*mongo.DeleteResult, error) {
	result, err := collection.DeleteOne(context.TODO(), bson.D{{Key: filter.Querykey, Value: filter.Queryvalue}})
	return result, err
}


