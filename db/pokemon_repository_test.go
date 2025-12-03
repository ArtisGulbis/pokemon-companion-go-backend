package db

import (
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models"
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
		ID:     25,
		Name:   "pikachu",
		Height: 4,
		Weight: 60,
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
}
