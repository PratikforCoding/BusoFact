package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "businfo"
const colName = "buslist"

var collection *mongo.Collection

func init() {
	err := godotenv.Load("./controllers/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connectString := os.Getenv("connectlink")
	clientOption := options.Client().ApplyURI(connectString)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	fmt.Println("MongoDB connection success")

	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection instance is ready!")

}

func GetBuses(source string, destination string) []primitive.M {
	filter := bson.M{
		"$and": []bson.M{
			{
				"stopages": bson.M{
					"$elemMatch": bson.M{
						"stopageName": source, // Source stopage name
					},
				},
			},
			{
				"stopages": bson.M{
					"$elemMatch": bson.M{
						"stopageName": destination, // Destination stopage name
					},
				},
			},
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var buses []primitive.M
	for cursor.Next(context.Background()) {
		var bus bson.M
		err := cursor.Decode(&bus)
		if err != nil {
			log.Fatal(err)
		}
		buses = append(buses, bus)
	}
	return buses
}

func AddBuses(name, stopageName string) (bson.M, error) {
	foundBus, err := GetBusByName(name)
	if err != nil {
		bus := bson.M{
			"name": name,
			"stopages": []bson.M{
				{"stopageNumber": 1, "stopage": stopageName},
			},
		}
		inserted, err := collection.InsertOne(context.Background(), bus)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted bus id:", inserted.InsertedID)
		updatedBus, err := GetBusByName(name)
		return updatedBus, nil
	}

	stopagesNum, ok := foundBus["stopages"].([]interface{})
    if !ok {
        fmt.Println("Invalid document format. 'stopages' field is missing or has an incorrect format.")
        return bson.M{}, errors.New("Wrong document format")
    }

	stopagesCount := len(stopagesNum) + 1

	filter := bson.M{"name": name}
	newStopage := bson.M{
        "stopageNumber": stopagesCount,
        "stopage":       stopageName,
    }
	update := bson.M{
        "$push": bson.M{"stopages": newStopage},
    }
	result, err := collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }
	if result.ModifiedCount == 1 {
        fmt.Println("Stopage added successfully")
    } else {
        fmt.Println("Stopage not added. Bus not found.")
    }
	updatedBus, err := GetBusByName(name)
	return updatedBus, nil
}

func GetBusByName(name string) (bson.M, error) {
	filter := bson.M{"name": name}
	var bus bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&bus)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("Bus not found")
			return bus , errors.New("Bus not found")
		} else {
			log.Fatal(err)
		}
	}
	
	return bus, nil
}

