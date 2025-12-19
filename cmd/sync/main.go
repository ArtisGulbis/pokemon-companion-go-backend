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
	flag.Parse()

	// 1. Setup database
	os.Remove("pokemon.db")
	database, err := db.New("pokemon.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// 2. Create API client and rate limiter (shared by all syncers)
	client := pokeapi.NewClient("https://pokeapi.co")
	rateLimiter := time.NewTicker(650 * time.Millisecond)
	defer rateLimiter.Stop()

	// 3. Create repositories
	versionRepo := db.NewVersionRepository(database)
	pokedexRepo := db.NewPokedexRepository(database)
	pokemonRepo := db.NewPokemonRepository(database)

	// 4. Create syncers (the building blocks)
	versionSyncer := services.NewVersionSyncer(client, versionRepo, rateLimiter)
	pokedexSyncer := services.NewPokedexSyncer(client, pokedexRepo, rateLimiter)
	pokemonSyncer := services.NewPokemonSyncer(client, pokemonRepo, rateLimiter)

	// 5. Create game syncer (the orchestrator)
	gameSyncer := services.NewGameSyncer(
		versionSyncer,
		pokedexSyncer,
		pokemonSyncer,
		rateLimiter,
	)

	// 6. Run the sync
	startTime := time.Now()

	if err := gameSyncer.SyncAllGames(1); err != nil {
		log.Fatal(err)
	}

	log.Printf("Sync complete! Time taken: %v", time.Since(startTime))
}
