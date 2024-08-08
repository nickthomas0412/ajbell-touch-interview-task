package main

import (
	"log"
	"net/http"
	"touch-coding-challenge/src/database"
	"touch-coding-challenge/src/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connection
	var err error
	db, err := database.InitDB()

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	err = database.SetupSchema(db)
	if err != nil {
		log.Fatalf("Error setting up schema: %v", err)
	}

	// Initialize HTTP routes exposing the db to the routes
	handlers.SetDatabase(db)
	setupRoutes()
}

// setupRoutes creates routes that are supported and launches the handler
func setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/clients", handlers.CreateClient).Methods("POST")
	router.HandleFunc("/deposits", handlers.CreateDeposit).Methods("POST")
	router.HandleFunc("/deposits/{id}", handlers.RetrieveDeposit).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
