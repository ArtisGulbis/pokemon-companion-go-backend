package db

import (
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models"
)

func TestInsertPokedex(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPokedexRepository(db)

	pokedex := &models.Pokedex{
		ID:           1,
		Name:         "Unova",
		IsMainSeries: true,
	}

	err = repo.InsertPokedex(pokedex)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	retrieved, err := repo.GetPokedexByID(pokedex.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve: %v", err)
	}
	if retrieved.Name != pokedex.Name {
		t.Errorf("Expected %s, go %s", pokedex.Name, retrieved.Name)
	}
}
