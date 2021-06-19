package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	User struct {
		ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" `
		Firstname string             `json:"firstname" bson:"firstname" validate:"required"`
		Lastname  string             `json:"lastname" bson:"lastname" validate:"required"`
		Email     string             `json:"email" bson:"email" validate:"required"`
		Password  string             `json:"password" bson:"password" validate:"required"`
		Token     string             `json:"token,omitempty" bson:"token,omitempty"`
	}

	LoginCredentials struct {
		Email    string `json:"email" bson:"email" validate:"required"`
		Password string `json:"password" bson:"password" validate:"required"`
	}
)

func (u User) MarshalJSON() ([]byte, error) {
	type user User // prevent recursion
	x := user(u)
	x.Password = ""
	return json.Marshal(x)
}
