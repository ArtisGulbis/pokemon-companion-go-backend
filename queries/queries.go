// Package queries contains all SQL queries used by the application.
// Queries are embedded at compile time using go:embed directives.
package queries

import (
	_ "embed" // Import for side effects to enable go:embed
)

//go:embed sql/version/version.sql
var InsertVersion string

//go:embed sql/pokemon/ability.sql
var InsertAbility string

//go:embed sql/version/version_group.sql
var InsertVersionGroup string

//go:embed sql/pokedex/version_group_pokedex.sql
var InsertVersionGroupPokedex string

//go:embed sql/pokedex/pokedex_entry.sql
var InsertPokedexEntry string

//go:embed sql/pokemon/pokemon.sql
var InsertPokemon string

//go:embed sql/pokemon/species.sql
var InsertSpecies string

//go:embed sql/types.sql
var InsertPokemonType string

//go:embed sql/pokemon/get_pokemon.sql
var GetPokemonByID string

//go:embed sql/pokedex/pokedex.sql
var InsertPokedex string

//go:embed sql/pokemon/types.sql
var InsertType string

//go:embed sql/pokedex/pokedex_description.sql
var InsertPokemonDescriptions string

//go:embed sql/pokedex/pokedex_pokemon_entry.sql
var InsertPokedexPokemonEntry string

//go:embed sql/pokedex/get_pokedex_pokemon_entry.sql
var GetPokedexPokemonEntryByID string

//go:embed sql/pokedex/get_pokedex_description.sql
var GetPokedexDescriptionsByPokedexID string

//go:embed sql/pokedex/get_pokedex.sql
var GetPokedexByID string

//go:embed sql/version/get_version.sql
var GetVersionByID string
