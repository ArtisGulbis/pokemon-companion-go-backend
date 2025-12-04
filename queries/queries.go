// Package queries contains all SQL queries used by the application.
// Queries are embedded at compile time using go:embed directives.
package queries

import (
	_ "embed" // Import for side effects to enable go:embed
)

// Pokemon queries

//go:embed sql/pokemon.sql
var InsertPokemon string

//go:embed sql/types.sql
var InsertPokemonType string

//go:embed sql/get_pokemon.sql
var GetPokemonByID string

//go:embed sql/pokedex.sql
var InsertPokedex string

//go:embed sql/get_pokedex.sql
var GetPokedexById string
