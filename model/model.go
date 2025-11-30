package model

type Pokemon struct {
	ID     int    `json:"id"` // Struct tags map JSON fields
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	//Moves  []PokemonMove `json:"hello"`
	Types []PokemonType `json:"types"` // Nested array
	Stats []PokemonStat `json:"stats"` // Nested array
}

type PokemonStat struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"` // Nested object
	Stat     Stat `json:"stat"`   // Nested object
}

type Stat struct {
	Name string `json:"name"`
}

type PokemonMove struct {
	Slot int      `json:"slot"`
	Type TypeInfo `json:"type"` // Nested object
}

type PokemonType struct {
	Slot int      `json:"slot"`
	Type TypeInfo `json:"type"` // Nested object
}

type TypeInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
