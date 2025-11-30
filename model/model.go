package model

type Pokemon struct {
	ID     int           `json:"id"`
	Name   string        `json:"name"`
	Height int           `json:"height"`
	Weight int           `json:"weight"`
	Types  []PokemonType `json:"types"`
}

type PokemonType struct {
	Slot int    `json:"slot"`
	Name string `json:"name"`
}
