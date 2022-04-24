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

func CreateContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	return ctx, cancel
}

func CreateConnectionToMongoDB() {
	ctx, cancel := CreateContext()
	defer cancel()
	MongoClient, _ = mongo.Connect(ctx, options.Client().ApplyURI(URI))
}

func CreateNewMovementEntry(roomName, text string) {
	collection := MongoClient.Database("movement").Collection(roomName)
	ctx, cancel := CreateContext()
	defer cancel()
	res, _ := collection.InsertOne(ctx, &MovementCache{
		Text: text,
	})

	fmt.Printf("Inserted %v into collection!\n", res.InsertedID)
}

func DeleteMovementCollection(roomName string) {
	collection := MongoClient.Database("movement").Collection(roomName)
	ctx, cancel := CreateContext()
	defer cancel()
	res, _ := collection.DeleteMany(
		ctx,
		bson.M{},
	)
	fmt.Printf("Deleted %v documents\n", res.DeletedCount)
}

func CheckIfMovementCacheExist(roomName string) bool {
	collection := MongoClient.Database("movement").Collection(roomName)
	ctx, cancel := CreateContext()
	defer cancel()
	count, _ := collection.CountDocuments(ctx, bson.M{})
	return count > 0
}

func GetMovement(roomName string) []MovementCache {
	collection := MongoClient.Database("movement").Collection(roomName)
	ctx, cancel := CreateContext()
	defer cancel()
	filterCursor, _ := collection.Find(ctx, bson.M{})
	defer filterCursor.Close(ctx)

	var movements []MovementCache
	filterCursor.All(ctx, &movements)

	return movements
}

func CloseConnectionMongoDB() {
	ctx, cancel := CreateContext()
	defer cancel()
	if MongoClient != nil {
		MongoClient.Disconnect(ctx)
	}
}
