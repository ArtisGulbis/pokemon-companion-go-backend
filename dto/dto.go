package dto

import "github.com/ArtisGulbis/pokemon-companion-go-backend/model"

type PokemonResponse struct {
	ID     int    `json:"id"` // Struct tags map JSON fields
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	//Moves  []PokemonMove `json:"hello"`
	Types []PokemonTypeResponse `json:"types"` // Nested array
	Stats []PokemonStat         `json:"stats"` // Nested array
}

type PokemonStat struct {
	BaseStat int          `json:"base_stat"`
	Effort   int          `json:"effort"` // Nested object
	Stat     StatResponse `json:"stat"`   // Nested object
}

type StatResponse struct {
	Name string `json:"name"`
}

type PokemonMoveResponse struct {
	Slot int              `json:"slot"`
	Type TypeInfoResponse `json:"type"` // Nested object
}

type PokemonTypeResponse struct {
	Slot int              `json:"slot"`
	Type TypeInfoResponse `json:"type"` // Nested object
}

type TypeInfoResponse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func ToPokemon(p *PokemonResponse) *model.Pokemon {
	resp := &model.Pokemon{
		ID:     p.ID,
		Name:   p.Name,
		Height: p.Height,
		Weight: p.Weight,
		Types:  make([]model.PokemonType, len(p.Types)),
	}

	for i, t := range p.Types {
		resp.Types[i] = model.PokemonType{
			Slot: t.Slot,
			Name: t.Type.Name,
		}
	}

	return resp
}
