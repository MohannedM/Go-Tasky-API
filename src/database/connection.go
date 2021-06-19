package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createClient() *mongo.Client {
	dbPass, _ := os.LookupEnv("DB_PASSWORD")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	connect := fmt.Sprintf("mongodb+srv://mohannedm:%s@cluster0-usvsi.mongodb.net/", dbPass)
	connection := options.Client().ApplyURI(connect)
	client, err := mongo.Connect(ctx, connection)
	if err != nil {
		log.Panic(err)
	}
	return client
}

func GetDatabase() *mongo.Database {
	return createClient().Database("gotasky")
}
