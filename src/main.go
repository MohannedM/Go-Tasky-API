package main

import (
	"TaskyBE/src/controllers"
	"TaskyBE/src/middlewares"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	//Environment Configuration
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	port, _ := os.LookupEnv("PORT")

	router := mux.NewRouter()

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      router,
	}

	router.HandleFunc("/users/register", controllers.Register).Methods("POST")
	router.HandleFunc("/users/login", controllers.Login).Methods("POST")
	tasksRouter := router.PathPrefix("/tasks").Subrouter()
	tasksRouter.Use(middlewares.AuthMiddleware)
	tasksRouter.HandleFunc("", controllers.CreateTask).Methods("POST")
	tasksRouter.HandleFunc("", controllers.GetTasks).Methods("GET")
	tasksRouter.HandleFunc("/assigned", controllers.GetAssignedTasks).Methods("GET")
	tasksRouter.HandleFunc("/created", controllers.GetCreatedTasks).Methods("GET")
	tasksRouter.HandleFunc("/{id}", controllers.GetTask).Methods("GET")
	tasksRouter.HandleFunc("/{id}", controllers.UpdateTask).Methods("PUT")
	tasksRouter.HandleFunc("/{id}", controllers.DeleteTask).Methods("DELETE")

	fmt.Println("Server running on port:", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
