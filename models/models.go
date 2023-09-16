package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bus struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty`
	Name string `json:"name,omitempty"`
	Stopages []Stopage `json:"stopages,omitempty"`
}

type Stopage struct {
	StopageNumber int `json:"stopageNumber"`
	StopageName string `json:"stopageName"`
}