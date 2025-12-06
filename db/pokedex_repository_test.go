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

func TestInsertPokedexPokemonEntry(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPokedexRepository(db)

	// First, create the Pokedex that the entries will reference
	pokedex := &models.Pokedex{
		ID:           1,
		Name:         "Kanto",
		IsMainSeries: true,
	}
	err = repo.InsertPokedex(pokedex)
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
