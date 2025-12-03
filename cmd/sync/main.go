package main

import (
	"flag"
	"log"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/db"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/pokeapi"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/services"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	limit := flag.Int("limit", 10, "Number of Pokemon to sync")
	flag.Parse()

	// 1. Setup database
	database, err := db.New("pokemon.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// 2. Create dependencies
	client := pokeapi.NewClient("https://pokeapi.co")
	repo := db.NewPokemonRepository(database)
	syncer := services.NewPokemonSyncer(client, repo)

	// 3. Sync Pokemon
	log.Printf("Syncing %d Pokemon...", *limit)
	if err := syncer.SyncAll(*limit); err != nil {
		log.Fatal(err)
	}

	log.Println("Sync complete!")
}
