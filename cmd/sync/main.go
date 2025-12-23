package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/db"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/igdb"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/pokeapi"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/services"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	flag.Parse()

	// 1. Setup database
	database, err := db.New("pokemon.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// 2. Create API client and rate limiter (shared by all syncers)
	client := pokeapi.NewClient("https://pokeapi.co")
	rateLimiter := time.NewTicker(650 * time.Millisecond)
	defer rateLimiter.Stop()

	// 3. Create IGDB client for game covers
	igdbClientID := os.Getenv("IGDB_CLIENT_ID")
	igdbClientSecret := os.Getenv("IGDB_CLIENT_SECRET")
	igdbClient := igdb.NewIGDBClient(igdbClientID, igdbClientSecret)

	// 4. Create repositories
	versionRepo := db.NewVersionRepository(database)
	pokedexRepo := db.NewPokedexRepository(database)
	pokemonRepo := db.NewPokemonRepository(database)

	// 5. Create syncers (the building blocks)
	versionSyncer := services.NewVersionSyncer(client, igdbClient, versionRepo, rateLimiter)
	pokedexSyncer := services.NewPokedexSyncer(client, pokedexRepo, rateLimiter)
	pokemonSyncer := services.NewPokemonSyncer(client, pokemonRepo, rateLimiter)

	// 6. Create game syncer (the orchestrator)
	gameSyncer := services.NewGameSyncer(
		versionSyncer,
		pokedexSyncer,
		pokemonSyncer,
		rateLimiter,
	)

	// 7. Run the sync
	startTime := time.Now()

	if err := gameSyncer.SyncAllGames(8); err != nil {
		log.Fatal(err)
	}

	log.Printf("Sync complete! Time taken: %v", time.Since(startTime))
}
