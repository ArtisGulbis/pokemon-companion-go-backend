package dto

import (
	"log"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/utils"
)

type Pokedex struct {
	ID           int                  `json:"id"`
	Name         string               `json:"name"`
	IsMainSeries bool                 `json:"isMainSeries"`
	Descriptions []PokedexDescription `json:"descriptions"`
	Pokemon      []PokemonEntry       `json:"pokemon"`
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

func NewPokedex(
	ext *external.Pokedex,
	descriptions []*external.PokedexDescriptions,
	entries []*external.PokedexPokemonEntry,
) *Pokedex {
	return &Pokedex{
		ID:           ext.ID,
		Name:         ext.Name,
		IsMainSeries: ext.IsMainSeries,
		Descriptions: mapDescriptions(descriptions),
		Pokemon:      mapPokemonEntries(entries),
	}
}

func mapDescriptions(descriptions []*external.PokedexDescriptions) []PokedexDescription {
	dtos := make([]PokedexDescription, len(descriptions))
	for i, d := range descriptions {
		dtos[i] = PokedexDescription{
			Language:    d.Language.Name,
			Description: d.Description,
		}
	}
	return dtos
}

func mapPokemonEntries(pokemonEntries []*external.PokedexPokemonEntry) []PokemonEntry {
	dtos := make([]PokemonEntry, len(pokemonEntries))
	for i, d := range pokemonEntries {
		id, err := utils.ExtractIDFromURL(d.PokemonSpecies.Url)
		if err != nil {
			log.Fatal(err)
		}
		dtos[i] = PokemonEntry{
			EntryNumber: d.EntryNumber,
			Name:        d.PokemonSpecies.Name,
			SpeciesID:   id,
		}
	}
	return dtos
}
