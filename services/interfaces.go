package services

import (
	"github.com/ArtisGulbis/pokemon-companion-go-backend/igdb"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
)

type PokemonAPIClient interface {
	FetchAll(path string) ([]external.Response, error)
	FetchPokemon(id int) (*external.Pokemon, error)
	FetchSpecies(id int) (*external.Species, error)
}

type VersionAPIClient interface {
	FetchAll(path string) ([]external.Response, error)
	FetchVersion(id int) (*external.Version, error)
	FetchVersionGroup(id int) (*external.VersionGroup, error)
}

type MoveAPIClient interface {
	FetchAll(path string) ([]external.Response, error)
	FetchMove(id int) (*external.Move, error)
}

type PokedexAPIClient interface {
	FetchAll(path string) ([]external.Response, error)
	FetchPokedex(id int) (*external.Pokedex, error)
}

type MoveRepo interface {
	InsertMove(v *external.Move) error
	GetMoveByID(id int) (*dto.Move, error)
}

type VersionRepo interface {
	InsertVersion(v *external.Version) error
	InsertVersionGroup(v *external.VersionGroup) error
	GetVersionByID(id int) (*dto.Version, error)
}

type PokemonRepo interface {
	InsertPokemon(p *external.Pokemon) error
	InsertType(p *external.PokemonType, pokemonId int) error
	InsertAbility(p *external.Ability, pokemonId int) error
	InsertSpecies(p *external.Species) error
	GetPokemonByID(id int) (*dto.Pokemon, error)
}

type PokedexRepo interface {
	InsertPokedex(p *external.Pokedex) error
	InsertPokedexEntry(p *external.PokedexEntry) error
	InsertVersionGroupPokedex(versionGroupPokedex *external.VersionGroup) error
	GetPokedexByID(id int) (*dto.Pokedex, error)
}

type IGDBClient interface {
	GetPokemonGameCover(versionName string) (*igdb.Game, error)
}
