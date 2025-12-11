package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/db"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/pokeapi"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/services"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	limit := flag.Int("limit", 10, "Number of Pokemon to sync")
	flag.Parse()

	// 1. Setup database
	os.Remove("pokemon.db")
	database, err := db.New("pokemon.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// 2. Create dependencies
	client := pokeapi.NewClient("https://pokeapi.co")
	rateLimiter := time.NewTicker(650 * time.Millisecond)
	defer rateLimiter.Stop()

	// pokemonRepo := db.NewPokemonRepository(database)
	//pokedexRepo := db.NewPokedexRepository(database)
	versionRepo := db.NewVersionRepository(database)

	//pokemonSyncer := services.NewPokemonSyncer(client, pokemonRepo, rateLimiter)
	versionSyncer := services.NewVersionSyncer(client, versionRepo, rateLimiter)

	startTime := time.Now()

	log.Printf("Syncing %d Versions...", *limit)
	if err := versionSyncer.SyncAll(*limit); err != nil {
		log.Fatal(err)
	}

	log.Printf("Sync complete! Time taken: %v", time.Since(startTime))
}
