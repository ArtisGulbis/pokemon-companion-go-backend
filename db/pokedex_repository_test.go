package db

import (
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/require"
)

func TestInsertPokedex(t *testing.T) {
	db := setupTest(t)

	pokedex := &external.Pokedex{
		ID:   1,
		Name: "Unova",
		Region: external.Response{
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

func TestInsertPokedexEntries(t *testing.T) {
	db := setupTest(t)

	pokedexEntry := &external.PokedexEntry{
		PokedexID:   1,
		SpeciesID:   1,
		EntryNumber: 1,
	}

	repo := NewPokedexRepository(db)

	_, err := db.Exec(`INSERT INTO species (id, name) VALUES (1, 'pikachu')`)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}
	_, err = db.Exec(`INSERT INTO pokedexes (id, name, region_name) VALUES (1, 'red', 'kanto')`)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	err = repo.InsertPokedexEntry(pokedexEntry)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}
	require.NoError(t, err)
}
