package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/db"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/middleware"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/routes"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	// Initialize the database connection
	database, err := db.New("pokemon.db")
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	router := routes.New(database)

	handler := router.Setup()

	// Create the middleware stack
	// This wraps our routes with logging and authorization
	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.Authorization,
	)

	// Configure the HTTP server
	server := http.Server{
		Addr:    ":3000",
		Handler: stack(handler), // Wrap the router with middleware
	}

	// Start the server
	fmt.Println("Server listening on port 3000")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
