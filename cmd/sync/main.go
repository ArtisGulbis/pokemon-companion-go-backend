package main

import (
	"fmt"
	"os"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/db"
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

}
