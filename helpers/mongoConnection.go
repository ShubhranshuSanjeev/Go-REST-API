package helpers

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DBClient MongoDB client
var DBClient *mongo.Client

//ConnectDB function to connect to MongoDB
func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	DBClient, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = DBClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}
