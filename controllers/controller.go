package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
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
		bus := model.Bus{
			Name: name,
			Stopages: []model.Stopage{
				{
					StopageNumber: 1,	
					StopageName: stopageName,
				},
			},
		}
		inserted, err := apiCfg.BusCollection.InsertOne(context.Background(), bus)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		fmt.Println("Inserted bus id:", inserted.InsertedID)
		NewBus, err := apiCfg.getBusByName(name)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		
		return NewBus, nil
	}
	return foundBus, nil
}

func (apiCfg *APIConfig) addBusStopage(name, stopage, beforeStopage string) (bson.M, error) {
    foundBus, err := apiCfg.getBusByName(name)
    if err != nil {
        return bson.M{}, errors.New("Couldn't find the bus")
    }

    // Get the existing stopages
    existingStopages := foundBus["stopages"].(primitive.A)

    // Find the existing stopage with "beforeStopage" and its position
    var beforeStopageIndex int
    var beforeStopageNumber int32

    for i, s := range existingStopages {
        stop := s.(primitive.M)
        if stop["stopage"].(string) == beforeStopage {
            beforeStopageIndex = i
            beforeStopageNumber = stop["stopageNumber"].(int32)
            break
        }
    }

    if beforeStopageIndex < 0 {
        return bson.M{}, errors.New("beforeStopage not found")
    }

    newStopageNumber := beforeStopageNumber + 1

    newStopage := bson.M{
        "stopageNumber": newStopageNumber,
        "stopage":      stopage,
    }

    // Create an empty primitive.A slice for updated stopages
    updatedStopages := make(primitive.A, 0, len(existingStopages)+1)

    for i, s := range existingStopages {
        stop := s.(primitive.M)
        updatedStopages = append(updatedStopages, stop)

        // Append the new stopage after the beforeStopage
        if i == beforeStopageIndex {
            updatedStopages = append(updatedStopages, newStopage)
        }

        // Increment the stopageNumber of stopages after the new one
        if i > beforeStopageIndex {
            stop["stopageNumber"] = stop["stopageNumber"].(int32) + 1
        }
    }

    // Update the bus document with the new stopages array
    filter := bson.M{"name": name}
    update := bson.M{
        "$set": bson.M{"stopages": updatedStopages},
    }

    _, err = apiCfg.BusCollection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
        return bson.M{}, err
    }

    updatedBus, err := apiCfg.getBusByName(name)
    if err != nil {
        log.Println(err)
        return bson.M{}, err
    }

    return updatedBus, nil
}

func (apiCfg *APIConfig) getBusByName(name string) (bson.M, error) {
    filter := bson.M{"name": name}
    var bus bson.M
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a reasonable timeout
    defer cancel()

    err := apiCfg.BusCollection.FindOne(ctx, filter).Decode(&bus)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            fmt.Println("Bus not found")
            return bus, errors.New("bus not found")
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

 

