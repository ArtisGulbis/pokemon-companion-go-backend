// Package queries contains all SQL queries used by the application.
// Queries are embedded at compile time using go:embed directives.
package queries

import (
	_ "embed" // Import for side effects to enable go:embed
)

// Pokemon queries

//go:embed pokemon.sql
var InsertPokemon string

//go:embed get_pokemon.sql
var GetPokemonByID string
