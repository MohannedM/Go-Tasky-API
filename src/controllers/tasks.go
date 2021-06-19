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

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const tasksTable = "tasks"

func taskValidator(taskStruct models.Task, res http.ResponseWriter) error {
	validate := validator.New()
	err := validate.Struct(taskStruct)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Println(validationErrors)
		ErrorThrower(http.StatusBadRequest, validationErrors.Error(), res, err)
		return err
	}
	return nil
}

func CreateTask(res http.ResponseWriter, req *http.Request) {
	var task models.Task
	var user models.User
	returnedUser := GetUser(&user, res, req)
	if returnedUser == nil {
		return
	}
	collection := database.GetDatabase().Collection(tasksTable)
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		ErrorThrower(http.StatusUnauthorized, fmt.Sprintf("%v", err), res, err)
		return
	}
	err = taskValidator(task, res)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	task.CreatedBy = user.ID
	i, err := collection.InsertOne(ctx, task)
	if err != nil {
		ErrorThrower(http.StatusUnauthorized, fmt.Sprintf("%v", err), res, err)
		return
	}
	task.ID = i.InsertedID

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	jsonData, _ := json.Marshal(task)
	res.Write(jsonData)
}

func GetTasks(res http.ResponseWriter, req *http.Request) {
	var user models.User
	returnedUser := GetUser(&user, res, req)
	if returnedUser == nil {
		return
	}
	collection := database.GetDatabase().Collection(tasksTable)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		ErrorThrower(http.StatusUnauthorized, fmt.Sprintf("%v", err), res, err)
		defer cursor.Close(ctx)
		return
	}

	results := []bson.M{}
	for cursor.Next(ctx) {
		var result bson.M
		cursor.Decode(&result)
		results = append(results, result)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	jsonData, _ := json.Marshal(results)
	res.Write(jsonData)
}

func GetAssignedTasks(res http.ResponseWriter, req *http.Request) {
	var user models.User
	returnedUser := GetUser(&user, res, req)
	if returnedUser == nil {
		return
	}
	collection := database.GetDatabase().Collection(tasksTable)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{"assignedTo": user.ID})
	if err != nil {
		ErrorThrower(http.StatusUnauthorized, fmt.Sprintf("%v", err), res, err)
		defer cursor.Close(ctx)
		return
	}

	results := []bson.M{}
	for cursor.Next(ctx) {
		var result bson.M
		cursor.Decode(&result)
		results = append(results, result)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	jsonData, _ := json.Marshal(results)
	res.Write(jsonData)
}

func GetCreatedTasks(res http.ResponseWriter, req *http.Request) {
	var user models.User
	returnedUser := GetUser(&user, res, req)
	if returnedUser == nil {
		return
	}
	collection := database.GetDatabase().Collection(tasksTable)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{"createdBy": user.ID})
	if err != nil {
		ErrorThrower(http.StatusUnauthorized, fmt.Sprintf("%v", err), res, err)
		defer cursor.Close(ctx)
		return
	}

	results := []bson.M{}
	for cursor.Next(ctx) {
		var result bson.M
		cursor.Decode(&result)
		results = append(results, result)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	jsonData, _ := json.Marshal(results)
	res.Write(jsonData)
}

func GetTask(res http.ResponseWriter, req *http.Request) {
	var task models.Task
	var user models.User
	returnedUser := GetUser(&user, res, req)
	if returnedUser == nil {
		return
	}
	params := mux.Vars(req)
	strId := params["id"]
	id, err := primitive.ObjectIDFromHex(strId)
	collection := database.GetDatabase().Collection(tasksTable)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	errDB := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)

	if err != nil {
		ErrorThrower(http.StatusUnauthorized, fmt.Sprintf("%v", err), res, err)
		return
	}

	if errDB != nil {
		ErrorThrower(http.StatusUnauthorized, "Could not find task", res, errDB)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	jsonData, _ := json.Marshal(task)
	res.Write(jsonData)
}

func UpdateTask(res http.ResponseWriter, req *http.Request) {
	var task models.Task
	var user models.User
	returnedUser := GetUser(&user, res, req)
	if returnedUser == nil {
		return
	}
	params := mux.Vars(req)
	strId := params["id"]
	id, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		ErrorThrower(http.StatusUnauthorized, fmt.Sprintf("%v", err), res, err)
		return
	}
	collection := database.GetDatabase().Collection(tasksTable)
	err = json.NewDecoder(req.Body).Decode(&task)

	if err != nil {
		ErrorThrower(http.StatusUnauthorized, fmt.Sprintf("%v", err), res, err)
		return
	}

	err = taskValidator(task, res)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	update := bson.M{
		"$set": task,
	}
	_, errDB := collection.UpdateByID(ctx, id, update)

	if errDB != nil {
		ErrorThrower(http.StatusUnauthorized, "Could not find task", res, errDB)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	jsonData, _ := json.Marshal(task)
	res.Write(jsonData)
}

func DeleteTask(res http.ResponseWriter, req *http.Request) {
	var task models.Task
	var user models.User
	returnedUser := GetUser(&user, res, req)
	if returnedUser == nil {
		return
	}
	params := mux.Vars(req)
	strId := params["id"]
	id, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		ErrorThrower(http.StatusUnauthorized, fmt.Sprintf("%v", err), res, err)
		return
	}
	collection := database.GetDatabase().Collection(tasksTable)
	err = json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		ErrorThrower(http.StatusUnauthorized, fmt.Sprintf("%v", err), res, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		ErrorThrower(http.StatusUnauthorized, "Could not find task", res, err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	jsonData, _ := json.Marshal(task)
	res.Write(jsonData)
}
