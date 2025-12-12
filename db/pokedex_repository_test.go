package db

import (
	"testing"

	models "github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
)

func TestInsertPokedex(t *testing.T) {
	db := setupTest(t)

	pokedex := &models.Pokedex{
		ID:   1,
		Name: "Unova",
		Region: models.Response{
			Name: "Unova",
			Url:  "some url",
		},
	}

	repo := NewPokedexRepository(db)

	err := repo.InsertPokedex(pokedex)
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
