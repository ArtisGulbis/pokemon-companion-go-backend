package services

import (
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
)

type PokemonAPIClient interface {
	FetchAll(path string) ([]external.Response, error)
	FetchPokemon(id int) (*external.Pokemon, error)
}

type VersionAPIClient interface {
	FetchAll(path string) ([]external.Response, error)
	FetchVersion(id int) (*external.Version, error)
	FetchVersionGroup(name string) (*external.VersionGroup, error)
}

type PokedexAPIClient interface {
	FetchAll(path string) ([]external.Response, error)
	FetchPokedex(id int) (*external.Pokedex, error)
}

type VersionRepo interface {
	InsertVersion(v *external.Version) error
	InsertVersionGroup(v *external.VersionGroup) error
	GetVersionByID(id int) (*dto.Version, error)
}

type PokemonRepo interface {
	InsertPokemon(p *external.Pokemon) error
	GetPokemonByID(id int) (*external.Pokemon, error)
}

type PokedexRepo interface {
	InsertPokedex(p *external.Pokedex) error
	GetPokedexByID(id int) (*external.Pokedex, error)
}
