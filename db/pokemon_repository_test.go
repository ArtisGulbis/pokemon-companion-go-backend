package db

import (
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertPokemon(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Insert the species first (required by foreign key)
	_, err = db.Exec(`INSERT INTO species (id, name) VALUES (1, 'pikachu')`)
	if err != nil {
		t.Fatalf("Failed to insert species: %v", err)
	}

	repo := NewPokemonRepository(db)

	// Test data
	pokemon := &external.Pokemon{
		ID:             1,
		Name:           "pikachu",
		Height:         1,
		Weight:         1,
		IsDefault:      true,
		BaseExperience: 1,
		SpeciesID:      1,
		Stats: []external.Stat{
			{
				BaseStat: 1,
				Stat: external.Response{
					Name: "hp",
					Url:  "",
				},
			},
			{
				BaseStat: 1,
				Stat: external.Response{
					Name: "attack",
					Url:  "",
				},
			},
			{
				BaseStat: 1,
				Stat: external.Response{
					Name: "defense",
					Url:  "",
				},
			},
			{
				BaseStat: 1,
				Stat: external.Response{
					Name: "special-attack",
					Url:  "",
				},
			},
			{
				BaseStat: 1,
				Stat: external.Response{
					Name: "special-defense",
					Url:  "",
				},
			},
			{
				BaseStat: 1,
				Stat: external.Response{
					Name: "speed",
					Url:  "",
				},
			},
		},
		Sprites: external.Sprite{
			Other: external.Other{
				OfficialArtwork: external.OfficialArtwork{
					FrontDefault: "",
					FrontShiny:   "",
				},
			},
		},
	}

	// Act
	err = repo.InsertPokemon(pokemon)

	// Assert
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	// Verify it was inserted
	actual, err := repo.GetPokemonByID(1)

	require.NoError(t, err)
	require.NotNil(t, actual)

	expected := &dto.Pokemon{
		ID:                 1,
		SpeciesID:          1,
		Name:               "pikachu",
		IsDefault:          true,
		Height:             1,
		Weight:             1,
		BaseExperience:     1,
		HP:                 1,
		Attack:             1,
		Defense:            1,
		SpecialAttack:      1,
		SpecialDefense:     1,
		Speed:              1,
		SpriteFrontDefault: "",
		SpriteFrontShiny:   "",
		SpriteArtwork:      "",
	}
	assert.Equal(t, expected, actual)
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

func TestGetStat(t *testing.T) {
	tests := []struct {
		name     string
		stats    []external.Stat
		key      string
		expected int
	}{
		{
			name: "finds stat at index 0 (hp)",
			stats: []external.Stat{
				{
					BaseStat: 45,
					Stat: external.Response{
						Name: "hp",
						Url:  "https://pokeapi.co/api/v2/stat/1/",
					},
				},
				{
					BaseStat: 49,
					Stat: external.Response{
						Name: "attack",
						Url:  "https://pokeapi.co/api/v2/stat/2/",
					},
				},
			},
			key:      "hp",
			expected: 45,
		},
		{
			name: "finds stat in middle of array (attack)",
			stats: []external.Stat{
				{
					BaseStat: 45,
					Stat: external.Response{
						Name: "hp",
						Url:  "https://pokeapi.co/api/v2/stat/1/",
					},
				},
				{
					BaseStat: 49,
					Stat: external.Response{
						Name: "attack",
						Url:  "https://pokeapi.co/api/v2/stat/2/",
					},
				},
				{
					BaseStat: 49,
					Stat: external.Response{
						Name: "defense",
						Url:  "https://pokeapi.co/api/v2/stat/3/",
					},
				},
			},
			key:      "attack",
			expected: 49,
		},
		{
			name: "finds special-attack with hyphen",
			stats: []external.Stat{
				{
					BaseStat: 65,
					Stat: external.Response{
						Name: "hp",
						Url:  "https://pokeapi.co/api/v2/stat/1/",
					},
				},
				{
					BaseStat: 55,
					Stat: external.Response{
						Name: "special-attack",
						Url:  "https://pokeapi.co/api/v2/stat/4/",
					},
				},
			},
			key:      "special-attack",
			expected: 55,
		},
		{
			name: "returns 0 when stat not found",
			stats: []external.Stat{
				{
					BaseStat: 45,
					Stat: external.Response{
						Name: "hp",
						Url:  "https://pokeapi.co/api/v2/stat/1/",
					},
				},
			},
			key:      "attack",
			expected: 0,
		},
		{
			name:     "returns 0 for empty stats array",
			stats:    []external.Stat{},
			key:      "hp",
			expected: 0,
		},
		{
			name:     "returns 0 for nil stats array",
			stats:    nil,
			key:      "hp",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getStat(tt.stats, tt.key)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestInsertType(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPokemonRepository(db)

	pokemonType := &external.PokemonType{
		Type: external.Response{
			Name: "grass",
			Url:  "https://pokeapi.co/api/v2/type/12/",
		},
		Slot: 2,
	}

	_, err = db.Exec(`INSERT INTO species (id, name) VALUES (1, "Pikachu")`)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}
	_, err = db.Exec(`INSERT INTO pokemon (id, species_id, name, is_default) VALUES (1, 1, "pikachu", FALSE)`)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}
	err = repo.InsertType(pokemonType, 1)
	require.NoError(t, err)
}

func TestInsertAbility(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPokemonRepository(db)

	ability := &external.Ability{
		IsHidden: false,
		Slot:     1,
		Ability: external.Response{
			Name: "bite",
			Url:  "https://pokeapi.co/api/v2/ability/39/",
		},
	}

	_, err = db.Exec(`INSERT INTO species (id, name) VALUES (1, "Pikachu")`)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}
	_, err = db.Exec(`INSERT INTO pokemon (id, species_id, name, is_default) VALUES (1, 1, "pikachu", FALSE)`)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}
	err = repo.InsertAbility(ability, 1)
	query := `SELECT * FROM pokemon_abilities WHERE pokemon_id = 1`
	var got dto.Ability
	err = db.QueryRow(query, 1).Scan(
		&got.PokemonId,
		&got.AbilityName,
		&got.IsHidden,
		&got.Slot,
	)
	require.NoError(t, err)

	expected := dto.Ability{
		PokemonId:   1,
		AbilityName: "bite",
		IsHidden:    false,
		Slot:        1,
	}

	assert.NotNil(t, got)
	assert.Equal(t, expected, got)
}
