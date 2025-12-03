package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type PokemonSyncer struct {
	client PokemonAPIClient
	repo   PokemonRepo
}

func NewPokemonSyncer(client PokemonAPIClient, repo PokemonRepo) *PokemonSyncer {
	return &PokemonSyncer{
		client: client,
		repo:   repo,
	}
}

// func (s *PokemonSyncer) SyncPokemon(id int) error {
// 	// 1. Fetch from API
// 	pokemon, err := s.client.FetchPokemon()
// 	if err != nil {
// 		return err
// 	}

// 	// 2. Save to database
// 	return s.repo.InsertPokemon(pokemon)
// }

// extractIDFromURL extracts the Pokemon ID from a URL like "https://pokeapi.co/api/v2/pokemon/25/"
func extractIDFromURL(url string) (int, error) {
	// Remove trailing slash and split by "/"
	url = strings.TrimSuffix(url, "/")
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid URL format: %s", url)
	}
	// Last part should be the ID
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse ID from URL %s: %w", url, err)
	}
	return id, nil
}

func (s *PokemonSyncer) SyncAll(limit int) error {
	// Fetch all Pokemon from API
	allPokemonResponse, err := s.client.FetchAll(fmt.Sprintf("pokemon?limit=%d", limit))
	if err != nil {
		log.Fatal(err)
	}
	for _, apr := range allPokemonResponse {
		// Extract ID from URL
		id, err := extractIDFromURL(apr.Url)
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
		fmt.Println("Inserted ", pokemon.Name)
	}
	return nil
}
