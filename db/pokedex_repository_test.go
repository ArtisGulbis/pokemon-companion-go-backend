package db

import (
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	models "github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGetPokedexByID_WithAllRelations(t *testing.T) {
	db := setupTest(t)
	repo := NewPokedexRepository(db)

	pokedex := &models.Pokedex{
		ID:   1,
		Name: "Kanto",
		Region: models.Response{
			Name: "Kanto",
			Url:  "",
		},
	}

	err := repo.InsertPokedex(pokedex)
	if err != nil {
		t.Fatalf("Failed to insert pokedex: %v", err)
	}

	got, err := repo.GetPokedexComplete(pokedex.ID)
	require.NoError(t, err)
	require.NotNil(t, got)

	expected := dto.Pokedex{
		ID:         1,
		Name:       "Kanto",
		RegionName: "Kanto",
	}

	assert.Equal(t, expected, *got)
}
