package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MovementCache struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Text string             `bson:"text"`
}

var (
	URI         string = os.Getenv("mongodb_uri")
	MongoClient *mongo.Client
)

func CreateContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

func CreateConnectionToMongoDB() {
	MongoClient, _ = mongo.Connect(CreateContext(), options.Client().ApplyURI(URI))
}

func CreateNewMovementEntry(roomName, text string) {
	collection := MongoClient.Database("movement").Collection(roomName)
	res, _ := collection.InsertOne(CreateContext(), &MovementCache{
		Text: text,
	})

	fmt.Printf("Inserted %v into collection!\n", res.InsertedID)
}

func DeleteMovementCollection(roomName string) {
	collection := MongoClient.Database("movement").Collection(roomName)
	res, _ := collection.DeleteMany(
		CreateContext(),
		bson.M{},
	)
	fmt.Printf("Deleted %v documents\n", res.DeletedCount)
}

func CheckIfMovementCacheExist(roomName string) bool {
	collection := MongoClient.Database("movement").Collection(roomName)
	count, _ := collection.CountDocuments(CreateContext(), bson.M{})
	return count > 0
}

func GetMovement(roomName string) []MovementCache {
	collection := MongoClient.Database("movement").Collection(roomName)
	filterCursor, _ := collection.Find(CreateContext(), bson.M{})
	defer filterCursor.Close(CreateContext())

	var movements []MovementCache
	filterCursor.All(CreateContext(), &movements)

	return movements
}

func CloseConnectionMongoDB() {
	if MongoClient != nil {
		MongoClient.Disconnect(CreateContext())
	}
}
