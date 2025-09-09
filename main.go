package main

import (
	"log"
	"net/http"

	"github.com/Sospeter-M-Kinyanjui/task-manager/database"
	"github.com/Sospeter-M-Kinyanjui/task-manager/handlers"
	"github.com/Sospeter-M-Kinyanjui/task-manager/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.Connect()

	r := mux.NewRouter()

	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.Auth)
	api.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	api.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
