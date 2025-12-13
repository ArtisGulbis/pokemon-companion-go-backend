package services

import (
	"fmt"
	"log"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/utils"
)

type PokemonSyncer struct {
	client      PokemonAPIClient
	repo        PokemonRepo
	rateLimiter *time.Ticker
}

func NewPokemonSyncer(client PokemonAPIClient, repo PokemonRepo, rateLimiter *time.Ticker) *PokemonSyncer {
	return &PokemonSyncer{
		client:      client,
		repo:        repo,
		rateLimiter: rateLimiter,
	}
}

func (s *PokemonSyncer) SyncSpecies(id int) error {
	return nil
}

// extractIDFromURL extracts the Pokemon ID from a URL like "https://pokeapi.co/api/v2/pokemon/25/"

func (s *PokemonSyncer) SyncAll(limit int) error {
	// Fetch all Pokemon from API
	allPokemonResponse, err := s.client.FetchAll(fmt.Sprintf("pokemon?limit=%d", limit))
	if err != nil {
		log.Fatal(err)
	}

	for i, apr := range allPokemonResponse {
		// Wait for rate limiter (except first request)
		if i > 0 {
			<-s.rateLimiter.C
		}

		// Extract ID from URL
		id, err := utils.ExtractIDFromURL(apr.Url)
		if err != nil {
			log.Fatal(err)
		}

		pokemon, err := s.client.FetchPokemon(id)
		if err != nil {
			log.Fatal(err)
		}
		err = s.repo.InsertPokemon(pokemon)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted Pokemon %s (%d/%d)\n", pokemon.Name, i+1, len(allPokemonResponse))
	}
	return nil
}
