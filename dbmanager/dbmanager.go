// Connects to MongoDB and sets a Stable API version
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type HashDocument struct {
	Hash string `bson:"hlf_hash"`
}

// Connecting to a MongoDB instance requres the usage of a MongoDB URI.
func connectMongo(db_uri string) *mongo.Client {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(db_uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)

	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}

func mongoSelectCollection(mongo_connection *mongo.Client, database_name string, collection_name string) *mongo.Collection {
	return mongo_connection.Database(database_name).Collection(collection_name)
}

func mongoUpdateHashBy_ID(mongo_collection *mongo.Collection, _id int, hash string) int {
	filter := bson.D{{"_id", _id}}
	update := bson.D{{"$set", bson.D{{"hlf_hash", hash}}}}
	result, err := mongo_collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Printf("mongoUpdateHashByID: Error %s\n", err)
		return 1
	}
	fmt.Printf("mongoUpdateHashByID: Successfully updated hash of %d\n", _id)
	return 0
}

func mongoUpdateHashByMatching(mongo_collection *mongo.Collection, filter bson.D, hash string) int {
	update := bson.D{{"$set", bson.D{{"hlf_hash", hash}}}}
	result, err := mongo_collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Printf("mongoUpdateHashByMatching: Error %s\n", err)
		return 1
	}
	fmt.Printf("mongoUpdateHashByMatching: Successfully updated hash of %d\n", result.UpsertedID)
	return 0
}

func mongoVerifyHashBy_ID(mongo_collection *mongo.Collection, _id int, hash string) bool {
	var result HashDocument
	filter := bson.D{{"_id", _id}}
	opts := options.FindOne().SetProjection(bson.D{{"hlf_hash", 1}})
	err := mongo_collection.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("mongoVerifyHashBy_ID: No documents found")
		} else {
			fmt.Println("mongoVerifyHashBy_ID: Error occured while verifying hash by ID.")
		}
	return false
	}
	return (result.Hash == hash)
}

func mongoVerifyHashByMatching(mongo_collection *mongo.Collection, filter bson.D, hash string) bool {
	var result HashDocument
	opts := options.FindOne().SetProjection(bson.D{{"hlf_hash", 1}})
	err := mongo_collection.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("mongoVerifyHashByMatching: No documents found")
		} else {
			fmt.Println("mongoVerifyHashByMatching: Error occured while verifying hash by ID.")
		}
		return false
	}
	return (result.Hash == hash)
}

func mongoDisconnect(mongo_connection *mongo.Client) int {
	err := mongo_connection.Disconnect(context.TODO())
	if err != nil {
		fmt.Println("mongoDisconnect: Failed to disconnect!")
		return 1
	}
	fmt.Println("mongoDisconnect: Successfully disconnected.")
	return 0
}