package db

import (
	"testing"

	models "github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
)

func TestInsertPokemon(t *testing.T) {
	// Create in-memory SQLite database for testing
	db, err := New(":memory:") // Special SQLite in-memory mode
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPokemonRepository(db)

	// Test data
	pokemon := &models.Pokemon{
		ID:             25,
		Name:           "pikachu",
		Height:         4,
		Weight:         60,
		BaseExperience: 25,
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
