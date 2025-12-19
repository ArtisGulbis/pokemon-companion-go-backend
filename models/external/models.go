package external

type Species struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	EvolutionChain URL      `json:"evolution_chain"`
	GenderRate     int      `json:"gender_rate"`
	CaptureRate    int      `json:"capture_rate"`
	BaseHappiness  int      `json:"base_happiness"`
	IsBaby         bool     `json:"is_baby"`
	IsLegendary    bool     `json:"is_legendary"`
	IsMythical     bool     `json:"is_mythical"`
	GrowthRate     Response `json:"growth_rate"`
	Generation     Response `json:"generation"`
}

type URL struct {
	URL string `json:"url"`
}

type Pokemon struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	Height         int           `json:"height"`
	Abilities      []Ability     `json:"abilities"`
	Weight         int           `json:"weight"`
	IsDefault      bool          `json:"is_default"`
	BaseExperience int           `json:"base_experience"`
	Types          []PokemonType `json:"types"`
	Stats          []Stat        `json:"stats"`
	Sprites        Sprite        `json:"sprites"`
	Species        Response      `json:"species"`
	SpeciesID      int
}

type Ability struct {
	IsHidden bool     `json:"is_hidden"`
	Ability  Response `json:"ability"`
	Slot     int      `json:"slot"`
}

type Sprite struct {
	Other Other `json:"other"`
}

type Other struct {
	OfficialArtwork OfficialArtwork `json:"official-artwork"`
}

type OfficialArtwork struct {
	FrontDefault string `json:"front_default"`
	FrontShiny   string `json:"front_shiny"`
}

type Stat struct {
	BaseStat int      `json:"base_stat"`
	Stat     Response `json:"stat"`
}

type Version struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Names        []Response `json:"names"`
	VersionGroup Response   `json:"version_group"`
}

type VersionGroup struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Generation Response   `json:"generation"`
	Pokedexes  []Response `json:"pokedexes"`
}

type Pokedex struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	Region         Response       `json:"region"`
	PokemonEntries []PokemonEntry `json:"pokemon_entries"`
}

type PokedexEntry struct {
	PokedexID   int
	SpeciesID   int
	EntryNumber int
}

type PokemonEntry struct {
	EntryNumber    int      `json:"entry_number"`
	PokemonSpecies Response `json:"pokemon_species"`
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
	Type Response `json:"type"`
	Slot int      `json:"slot"`
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
