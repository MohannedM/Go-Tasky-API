package controllers

import (
	"TaskyBE/src/database"
	"TaskyBE/src/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	gorillaContext "github.com/gorilla/context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Error struct {
	Message string
}

func ErrorThrower(status int, errMsg string, res http.ResponseWriter, err ...error) {
	log.Println(err)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	errObj := Error{errMsg}
	messageJson, _ := json.Marshal(errObj)
	res.Write([]byte(messageJson))
}

func GetUser(user *models.User, res http.ResponseWriter, req *http.Request) *models.User {
	userId := gorillaContext.Get(req, "user_id")
	userCollection := database.GetDatabase().Collection(userTable)
	userCtx, userCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer userCancel()
	strUserId := fmt.Sprintf("%v", userId)
	userId, err := primitive.ObjectIDFromHex(strUserId)
	if err != nil {
		ErrorThrower(http.StatusUnauthorized, "Could not find user", res, err)
		return nil
	}
	errDB := userCollection.FindOne(userCtx, bson.M{"_id": userId}).Decode(user)
	if errDB != nil {
		ErrorThrower(http.StatusUnauthorized, "Could not find user", res, errDB)
		return nil
	}
	return user
}
