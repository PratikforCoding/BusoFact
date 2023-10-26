package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bus struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string            `json:"name,omitempty" bson:"name,omitempty"`
	Stopages []Stopage         `json:"stopages,omitempty" bson:"stopages,omitempty"`
}

type Stopage struct {
	StopageNumber int    `json:"stopageNumber" bson:"stopageNumber"`
	StopageName   string `json:"stopageName" bson:"stopageName"`
}

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FristName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
	AccessToken string `json:"accesstoken"`
	RefreshToken string `json:"refreshToken"`
}