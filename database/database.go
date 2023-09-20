package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "businfo"
const busColName = "buses"
const usrColName = "users"

var busCollection *mongo.Collection
var usrCollection *mongo.Collection


var client *mongo.Client

func CreateDB(mongoUri string) (*mongo.Collection, *mongo.Collection, error) {
	clientOption := options.Client().ApplyURI(mongoUri)
	var err error
	client, err = mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	fmt.Println("MongoDB connection success")
	busCollection = client.Database(dbName).Collection(busColName)
	usrCollection = client.Database(dbName).Collection(usrColName)
	fmt.Println("Collection instance is ready")

	return busCollection, usrCollection, nil	
}



func CloseDB() {
	if client != nil {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("MongoDB connection closed")
	}
}