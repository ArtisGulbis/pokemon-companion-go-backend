package db

import (
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	models "github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) *PokedexRepository {
	db, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	repo := NewPokedexRepository(db)

	t.Cleanup(func() { db.Close() })

	return repo
}

func TestInsertPokedex(t *testing.T) {
	repo := setupTest(t)

	pokedex := &models.Pokedex{
		ID:           1,
		Name:         "Unova",
		IsMainSeries: true,
	}

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

func TestInsertPokedexDescription(t *testing.T) {
	repo := setupTest(t)

	pokedexDescriptions := []models.PokedexDescriptions{
		{
			Description: "description",
			Language: models.Response{
				Name: "de",
				Url:  "some url",
			},
		},
		{
			Description: "description2",
			Language: models.Response{
				Name: "en",
				Url:  "some url 2",
			},
		},
	}

	pokedex := &models.Pokedex{
		ID:           1,
		Name:         "Kanto",
		IsMainSeries: true,
	}

	err := repo.InsertPokedex(pokedex)
	if err != nil {
		t.Fatalf("Failed to insert pokedex: %v", err)
	}

	err = repo.InsertPokedexDescriptions(pokedexDescriptions, pokedex.ID)
	if err != nil {
		t.Fatalf("Failed to insert pokedex description: %v", err)
	}

	got, err := repo.GetPokedexDescriptionsByPokedexID(pokedex.ID)
	if err != nil {
		t.Fatalf("Failed to get pokedex description: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("Got description length %d: want %d", len(got), len(pokedexDescriptions))
	}

	for i, desc := range got {
		want := pokedexDescriptions[i]
		if desc.Description != want.Description {
			t.Fatalf("Got description: %s want: %s", desc.Description, want.Description)
		}
		if desc.Language.Name != want.Language.Name {
			t.Fatalf("Got language: %s want: %s", desc.Language.Name, want.Language.Name)
		}
	}

}

func TestInsertPokedexPokemonEntry(t *testing.T) {
	repo := setupTest(t)

	pokedex := &models.Pokedex{
		ID:           1,
		Name:         "Kanto",
		IsMainSeries: true,
	}

	err := repo.InsertPokedex(pokedex)
	if err != nil {
		t.Fatalf("Failed to insert pokedex: %v", err)
	}

	pokedexPokemonEntries := []models.PokedexPokemonEntry{
		{
			EntryNumber: 1,
			PokemonSpecies: models.Response{
				Name: "Bulbasaur",
				Url:  "https://pokeapi.co/api/v2/pokemon-species/1/",
			},
		},
		{
			EntryNumber: 2,
			PokemonSpecies: models.Response{
				Name: "Ivysaur",
				Url:  "https://pokeapi.co/api/v2/pokemon-species/2/",
			},
		},
	}

	err = repo.InsertPokedexPokemonEntry(pokedexPokemonEntries, 1)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	retrieved, err := repo.GetPokedexEntriesByPokedexID(1)
	if err != nil {
		t.Fatalf("Failed to retrieve: %v", err)
	}

	for i, e := range retrieved {
		currentEntry := pokedexPokemonEntries[i]
		if e.EntryNumber != currentEntry.EntryNumber {
			t.Errorf("Expected %d, go %d", currentEntry.EntryNumber, e.EntryNumber)
		}
		if e.PokemonSpecies.Name != currentEntry.PokemonSpecies.Name {
			t.Errorf("Expected %s, go %s", currentEntry.PokemonSpecies.Name, e.PokemonSpecies.Name)
		}
		if e.PokemonSpecies.Url != currentEntry.PokemonSpecies.Url {
			t.Errorf("Expected %s, go %s", currentEntry.PokemonSpecies.Url, e.PokemonSpecies.Url)
		}
	}

}

func TestGetPokedexByID_WithAllRelations(t *testing.T) {
	repo := setupTest(t)

	pokedex := &models.Pokedex{
		ID:           1,
		Name:         "Kanto",
		IsMainSeries: true,
	}

	err := repo.InsertPokedex(pokedex)
	if err != nil {
		t.Fatalf("Failed to insert pokedex: %v", err)
	}

	pokedexPokemonEntries := []models.PokedexPokemonEntry{
		{
			EntryNumber: 1,
			PokemonSpecies: models.Response{
				Name: "Bulbasaur",
				Url:  "https://pokeapi.co/api/v2/pokemon-species/1/",
			},
		},
		{
			EntryNumber: 2,
			PokemonSpecies: models.Response{
				Name: "Ivysaur",
				Url:  "https://pokeapi.co/api/v2/pokemon-species/2/",
			},
		},
	}

	err = repo.InsertPokedexPokemonEntry(pokedexPokemonEntries, 1)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	pokedexDescriptions := []models.PokedexDescriptions{
		{
			Description: "description",
			Language: models.Response{
				Name: "de",
				Url:  "some url",
			},
		},
		{
			Description: "description2",
			Language: models.Response{
				Name: "en",
				Url:  "some url 2",
			},
		},
	}

	err = repo.InsertPokedexDescriptions(pokedexDescriptions, pokedex.ID)
	if err != nil {
		t.Fatalf("Failed to insert pokedex description: %v", err)
	}

	got, err := repo.GetPokedexComplete(pokedex.ID)
	require.NoError(t, err)
	require.NotNil(t, got)

	expected := dto.Pokedex{
		ID:           1,
		Name:         "Kanto",
		IsMainSeries: true,
		Descriptions: []dto.PokedexDescription{
			{
				Language:    "de",
				Description: "description",
			},
			{
				Language:    "en",
				Description: "description2",
			},
		},
		Pokemon: []dto.PokemonEntry{
			{
				EntryNumber: 1,
				Name:        "Bulbasaur",
				SpeciesID:   1,
			},
			{
				EntryNumber: 2,
				Name:        "Ivysaur",
				SpeciesID:   2,
			},
		},
	}

	assert.Equal(t, expected, *got)
}
