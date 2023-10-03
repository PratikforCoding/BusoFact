package controller

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/PratikforCoding/BusoFact.git/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/PratikforCoding/BusoFact.git/models"
)

type APIConfig struct {
	jwtSecret string
	BusCollection *mongo.Collection
	UserCollection *mongo.Collection
}

func NewAPIConfig(busCol, usrCol *mongo.Collection, jwtSecret string) *APIConfig {
	return &APIConfig{BusCollection: busCol, UserCollection: usrCol, jwtSecret: jwtSecret}
}

func (apiCfg *APIConfig)getBuses(source string, destination string) []primitive.M {
	filter := bson.M{
		"$and": []bson.M{
			{
				"stopages": bson.M{
					"$elemMatch": bson.M{
						"stopage": source, // Source stopage name
					},
				},
			},
			{
				"stopages": bson.M{
					"$elemMatch": bson.M{
						"stopage": destination, // Destination stopage name
					},
				},
			},
		},
	}

	cursor, err := apiCfg.BusCollection.Find(context.Background(), filter)
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

func (apiCfg *APIConfig)addBuses(name, stopageName string) (bson.M, error) {
	foundBus, err := apiCfg.getBusByName(name)
	if err != nil {
		bus := bson.M{
			"name": name,
			"stopages": []bson.M{
				{"stopageNumber": 1, "stopage": stopageName},
			},
		}
		inserted, err := apiCfg.BusCollection.InsertOne(context.Background(), bus)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted bus id:", inserted.InsertedID)
		updatedBus, err := apiCfg.getBusByName(name)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return updatedBus, nil
	}

	stopagesNum, ok := foundBus["stopages"].(primitive.A)
    if !ok {
        fmt.Println("Invalid document format. 'stopages' field is missing or has an incorrect format.")
        return bson.M{}, errors.New("wrong document format")
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
	result, err := apiCfg.BusCollection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }
	if result.ModifiedCount == 1 {
        fmt.Println("Stopage added successfully")
    } else {
        fmt.Println("Stopage not added. Bus not found.")
    }
	updatedBus, err := apiCfg.getBusByName(name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return updatedBus, nil
}

func (apiCfg *APIConfig)getBusByName(name string) (bson.M, error) {
	filter := bson.M{"name": name}
	var bus bson.M
	err := apiCfg.BusCollection.FindOne(context.TODO(), filter).Decode(&bus)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("Bus not found")
			return bus , errors.New("bus not found")
		} else {
			log.Fatal(err)
		}
	}
	
	return bus, nil
}

func (apiCfg *APIConfig)createUser(firstName, lastName, email, password string) (model.User, error) {
	foundUser, err := apiCfg.getUser(email)
	if err != nil {
		hash, err := auth.HashedPassword(password)
		if err != nil {
			return model.User{}, err
		}
		user := model.User{
			FristName: firstName,
			LastName: lastName,
			Email: email,
			Password: hash,
		}

		inserted, err := apiCfg.UserCollection.InsertOne(context.Background(), user)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Inserted user id:", inserted.InsertedID)
		createdUser, err := apiCfg.getUser(email)
		if err != nil {
			log.Println(err)
			return model.User{}, err
		}
		return createdUser, nil
	}
	return foundUser, errors.New("user already exists")
}

func (apiCfg *APIConfig)userLogin(email, password string) (model.User, error) {
	var user model.User
	user, err := apiCfg.getUser(email)
	if err != nil {
		return model.User{}, errors.New("user doesn't exist")
	}
	userHash := user.Password
	err = auth.CheckPasswordHash(password, userHash)
	if err != nil {
		log.Println(err)
		return model.User{}, errors.New("wrong password")
	}
	return user, nil
}

func (apiCfg *APIConfig)getUser(email string) (model.User, error) {
	filter := bson.M{"email":email}
	var user model.User
	err := apiCfg.UserCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("user not found")
			return model.User{} , errors.New("user not found")
		} else {
			log.Fatal(err)
		}
	}
	
	return user, nil
}

 

