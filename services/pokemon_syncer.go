package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/utils"
)

type PokemonSyncer struct {
	client        PokemonAPIClient
	repo          PokemonRepo
	rateLimiter   *time.Ticker
	syncedSpecies map[int]bool // In-memory cache of synced species IDs
	syncedPokemon map[int]bool // In-memory cache of synced Pokemon IDs
	mu            sync.Mutex   // Protects cache maps
}

func NewPokemonSyncer(client PokemonAPIClient, repo PokemonRepo, rateLimiter *time.Ticker) *PokemonSyncer {
	return &PokemonSyncer{
		client:        client,
		repo:          repo,
		rateLimiter:   rateLimiter,
		syncedSpecies: make(map[int]bool),
		syncedPokemon: make(map[int]bool),
	}
}

func (s *PokemonSyncer) FetchSpecies(id int) (*external.Species, error) {
	return s.client.FetchSpecies(id)
}

func (s *PokemonSyncer) InsertSpecies(sp *external.Species) error {
	return s.repo.InsertSpecies(sp)
}

func (s *PokemonSyncer) InsertPokemon(p *external.Pokemon) error {
	return s.repo.InsertPokemon(p)
}

func (s *PokemonSyncer) InsertType(t *external.PokemonType, pokemonID int) error {
	return s.repo.InsertType(t, pokemonID)
}

func (s *PokemonSyncer) InsertAbility(a *external.Ability, pokemonID int) error {
	return s.repo.InsertAbility(a, pokemonID)
}

func (s *PokemonSyncer) SyncPokemon(id int) (*external.Pokemon, error) {
	// Check cache first
	s.mu.Lock()
	if s.syncedPokemon[id] {
		s.mu.Unlock()
		// Already synced in this session, just fetch from client without DB insert
		// We still return the Pokemon data for the caller
		pokemon, err := s.client.FetchPokemon(id)
		if err != nil {
			return nil, err
		}
		speciesID, err := utils.ExtractIDFromURL(pokemon.Species.Url)
		if err != nil {
			return nil, err
		}
		pokemon.SpeciesID = speciesID
		return pokemon, nil
	}
	s.mu.Unlock()

	// Not in cache, fetch and insert
	pokemon, err := s.client.FetchPokemon(id)
	if err != nil {
		return nil, err
	}

	speciesID, err := utils.ExtractIDFromURL(pokemon.Species.Url)
	if err != nil {
		return nil, err
	}
	pokemon.SpeciesID = speciesID

	// Download and save Pokemon sprites locally
	if pokemon.Sprites.Other.OfficialArtwork.FrontDefault != "" {
		localPath := utils.GetPokemonSpritePath(pokemon.ID, "artwork")
		if err := utils.DownloadImage(pokemon.Sprites.Other.OfficialArtwork.FrontDefault, localPath); err != nil {
			log.Printf("Warning: failed to download artwork for Pokemon %d: %v", pokemon.ID, err)
		} else {
			pokemon.Sprites.Other.OfficialArtwork.FrontDefault = localPath
		}
	}

	if pokemon.Sprites.Other.OfficialArtwork.FrontShiny != "" {
		localPath := utils.GetPokemonSpritePath(pokemon.ID, "artwork_shiny")
		if err := utils.DownloadImage(pokemon.Sprites.Other.OfficialArtwork.FrontShiny, localPath); err != nil {
			log.Printf("Warning: failed to download shiny artwork for Pokemon %d: %v", pokemon.ID, err)
		} else {
			pokemon.Sprites.Other.OfficialArtwork.FrontShiny = localPath
		}
	}

	err = s.repo.InsertPokemon(pokemon)
	if err != nil {
		return nil, err
	}

	// Mark as synced in cache
	s.mu.Lock()
	s.syncedPokemon[id] = true
	s.mu.Unlock()

	return pokemon, nil
}

func (s *PokemonSyncer) SyncSpecies(id int) (*external.Species, error) {
	// Check cache first
	s.mu.Lock()
	if s.syncedSpecies[id] {
		s.mu.Unlock()
		// Already synced in this session, skip API call and DB insert
		return nil, nil
	}
	s.mu.Unlock()

	// Not in cache, fetch and insert
	species, err := s.client.FetchSpecies(id)
	if err != nil {
		return nil, err
	}

	err = s.repo.InsertSpecies(species)
	if err != nil {
		return nil, err
	}

	// Mark as synced in cache
	s.mu.Lock()
	s.syncedSpecies[id] = true
	s.mu.Unlock()

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
