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
	"github.com/joho/godotenv"
)

func main() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.New("pokemon.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	client := pokeapi.NewClient("https://pokeapi.co")
	rateLimiter := time.NewTicker(650 * time.Millisecond)
	defer rateLimiter.Stop()

	igdbClientID := os.Getenv("IGDB_CLIENT_ID")
	igdbClientSecret := os.Getenv("IGDB_CLIENT_SECRET")
	igdbClient := igdb.NewIGDBClient(igdbClientID, igdbClientSecret)

	versionRepo := db.NewVersionRepository(database)
	pokedexRepo := db.NewPokedexRepository(database)
	pokemonRepo := db.NewPokemonRepository(database)
	moveRepo := db.NewMoveRepository(database)

	versionSyncer := services.NewVersionSyncer(client, igdbClient, versionRepo, rateLimiter)
	pokedexSyncer := services.NewPokedexSyncer(client, pokedexRepo, rateLimiter)
	pokemonSyncer := services.NewPokemonSyncer(client, pokemonRepo, rateLimiter)
	moveSyncer := services.NewMoveSyncer(client, moveRepo, rateLimiter)

	gameSyncer := services.NewGameSyncer(
		versionSyncer,
		pokedexSyncer,
		pokemonSyncer,
		moveSyncer,
		rateLimiter,
	)

	startTime := time.Now()

	if err := gameSyncer.SyncAllGames(100); err != nil {
		log.Fatal(err)
	}

	// scraper := scraper.NewScraper()
	// scraper.ScrapeGamePage("https://bulbapedia.bulbagarden.net/wiki/Pokémon_Gold_and_Silver_Versions")

	log.Printf("Sync complete! Time taken: %v", time.Since(startTime))
}
