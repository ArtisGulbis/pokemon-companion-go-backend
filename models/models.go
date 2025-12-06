package models

type Pokemon struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	BaseExperience int           `json:"base_experience"`
	Types          []PokemonType `json:"types"`
}

type Pokedex struct {
	ID             int                   `json:"id"`
	IsMainSeries   bool                  `json:"is_main_series"`
	Name           string                `json:"name"`
	Descriptions   []PokedexDescriptions `json:"descriptions"`
	PokemonEntries []PokedexPokemonEntry `json:"pokemon_entries"`
}

type PokedexPokemonEntry struct {
	EntryNumber    int      `json:"entry_number"`
	PokemonSpecies Response `json:"pokemon_species"`
}

type PokedexDescriptions struct {
	Description string   `json:"description"`
	Language    Response `json:"language"`
}

type PokemonType struct {
	Name string `json:"name"`
	Slot int    `json:"slot"`
}

type Response struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PaginatedResponse struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Response `json:"results"`
}
