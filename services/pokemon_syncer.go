package services

import (
	"fmt"
	"log"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
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

func (s *PokemonSyncer) FetchSpecies(id int) (*external.Species, error) {
	return s.client.FetchSpecies(id)
}

func (s *PokemonSyncer) InsertSpecies(sp *external.Species) error {
	return s.repo.InsertSpecies(sp)
}

func (s *PokemonSyncer) SyncSpecies(id int) (*external.Species, error) {
	species, err := s.client.FetchSpecies(id)
	if err != nil {
		return nil, err
	}

	err = s.repo.InsertSpecies(species)
	if err != nil {
		return nil, err
	}

	return species, nil
}

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
