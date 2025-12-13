package db

import (
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/require"
)

func TestInsertPokemon(t *testing.T) {
	// Create in-memory SQLite database for testing
	db, err := New(":memory:") // Special SQLite in-memory mode
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Insert the species first (required by foreign key)
	_, err = db.Exec(`INSERT INTO species (id, name, is_legendary, is_mythical) VALUES (25, 'pikachu', FALSE, FALSE)`)
	if err != nil {
		t.Fatalf("Failed to insert species: %v", err)
	}

	repo := NewPokemonRepository(db)

	// Test data
	pokemon := &external.Pokemon{
		ID:             25,
		Name:           "pikachu",
		Height:         4,
		Weight:         60,
		BaseExperience: 25,
		SpeciesID:      25,
		IsDefault:      true,
	}

	// Act
	err = repo.InsertPokemon(pokemon)

	// Assert
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	// Verify it was inserted
	retrieved, err := repo.GetPokemonByID(25)
	if err != nil {
		t.Fatalf("Failed to retrieve: %v", err)
	}

	if retrieved.Name != "pikachu" {
		t.Errorf("Expected pikachu, got %s", retrieved.Name)
	}

	if retrieved.BaseExperience != pokemon.BaseExperience {
		t.Errorf("Expected base_experience %d, got %d", pokemon.BaseExperience, retrieved.BaseExperience)
	}
	if retrieved.ID != pokemon.ID {
		t.Errorf("Expected id %d, got %d", pokemon.ID, retrieved.ID)
	}
	if retrieved.Weight != pokemon.Weight {
		t.Errorf("Expected weight %d, got %d", pokemon.Weight, retrieved.Weight)
	}
	if retrieved.Height != pokemon.Height {
		t.Errorf("Expected height %d, got %d", pokemon.Height, retrieved.Height)
	}
}

func TestInsertSpecies(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPokemonRepository(db)

	species := &external.Species{
		ID:            152,
		Name:          "chikorita",
		IsBaby:        false,
		IsLegendary:   false,
		IsMythical:    false,
		BaseHappiness: 70,
		CaptureRate:   45,
		EvolutionChain: external.URL{
			URL: "https://pokeapi.co/api/v2/evolution-chain/79/",
		},
		GrowthRate: external.Response{
			Name: "medium-slow",
			Url:  "https://pokeapi.co/api/v2/growth-rate/4/",
		},
		GenderRate: 1,
		Generation: external.Response{
			Name: "generation-ii",
			Url:  "https://pokeapi.co/api/v2/generation/2/",
		},
	}

	err = repo.InsertSpecies(species)
	require.NoError(t, err)
}

func TestInsertVersionGroupPokedex(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPokedexRepository(db)

	// Test data
	versionGroupPokedex := &external.VersionGroup{
		ID:   1,
		Name: "johto",
		Generation: external.Response{
			Name: "",
			Url:  "",
		},
		Pokedexes: []external.Response{
			{
				Name: "johto",
				Url:  "https://pokeapi.co/api/v2/pokedex/1/",
			},
		},
	}

	_, err = db.Exec(`INSERT INTO version_groups (id, name, generation_name) VALUES (1, "yellow", "yellow")`)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	_, err = db.Exec(`INSERT INTO pokedexes (id, name, region_name) VALUES (1, "yellow", "yellow")`)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	err = repo.InsertVersionGroupPokedex(versionGroupPokedex)
	require.NoError(t, err)
}
