package dto

type Pokedex struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	RegionName string `json:"regionName"`
}

type LocalizedName struct {
	Language    string `json:"language"`
	Description string `json:"description"`
}

type PokedexDescription struct {
	Language    string `json:"language"`
	Description string `json:"description"`
}

type PokemonEntry struct {
	EntryNumber int    `json:"entryNumber"`
	Name        string `json:"name"`
	SpeciesID   int    `json:"speciesId"`
}

type VersionGroupPokedex struct {
	VersionGroupID int
	PokedexID      int
}
