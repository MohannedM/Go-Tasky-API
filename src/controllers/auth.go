package controllers

import (
	"TaskyBE/src/database"
	"TaskyBE/src/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

const userTable = "users"

func CreateToken(userid interface{}) (string, error) {
	var err error
	secretKey, _ := os.LookupEnv("JWT_SECRET")
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func Register(res http.ResponseWriter, req *http.Request) {
	var user models.User
	collection := database.GetDatabase().Collection(userTable)
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		ErrorThrower(http.StatusBadRequest, fmt.Sprintf("%v", err), res, err)
		return
	}
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Println(validationErrors)
		ErrorThrower(http.StatusBadRequest, validationErrors.Error(), res, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	errDB := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&user)

	if errDB == nil {
		ErrorThrower(http.StatusUnauthorized, "Email already exists", res)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ErrorThrower(http.StatusInternalServerError, "Internal Error", res, err)
		return
	}
	user.Password = string(hashedPassword)

	i, err := collection.InsertOne(ctx, user)

	if err != nil {
		ErrorThrower(http.StatusInternalServerError, "An error occurred with database", res, err)
		return
	}

	token, _ := CreateToken(i.InsertedID)
	user.Token = token
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	jsonData, _ := user.MarshalJSON()
	res.Write(jsonData)
}

func Login(res http.ResponseWriter, req *http.Request) {
	var user models.User
	var loginCredentials models.LoginCredentials
	collection := database.GetDatabase().Collection(userTable)
	err := json.NewDecoder(req.Body).Decode(&loginCredentials)
	if err != nil {
		ErrorThrower(http.StatusBadRequest, fmt.Sprintf("%v", err), res, err)
		return
	}
	validate := validator.New()
	err = validate.Struct(loginCredentials)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Println(validationErrors)
		ErrorThrower(http.StatusBadRequest, validationErrors.Error(), res, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	errDB := collection.FindOne(ctx, bson.M{"email": loginCredentials.Email}).Decode(&user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCredentials.Password))

	if err != nil {
		ErrorThrower(http.StatusBadRequest, fmt.Sprintf("%v", err), res, err)
		return
	}

	if errDB != nil {
		ErrorThrower(http.StatusUnauthorized, "Email or password is incorrect", res, err)
		return
	}

	token, _ := CreateToken(user.ID)
	user.Token = token
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	jsonData, _ := user.MarshalJSON()
	res.Write(jsonData)
}
