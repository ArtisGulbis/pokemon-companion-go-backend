package services

import (
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	models "github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
)

type PokemonAPIClient interface {
	FetchAll(path string) ([]models.Response, error)
	FetchPokemon(id int) (*models.Pokemon, error)
}

type VersionAPIClient interface {
	FetchAll(path string) ([]models.Response, error)
	FetchVersion(id int) (*models.Version, error)
}

type PokedexAPIClient interface {
	FetchAll(path string) ([]models.Response, error)
	FetchPokedex(id int) (*models.Pokedex, error)
}

type VersionRepo interface {
	InsertVersion(v *external.Version) error
	GetVersionByID(id int) (*dto.Version, error)
}

type PokemonRepo interface {
	InsertPokemon(p *models.Pokemon) error
	GetPokemonByID(id int) (*models.Pokemon, error)
}

type PokedexRepo interface {
	InsertPokedex(p *models.Pokedex) error
	GetPokedexByID(id int) (*models.Pokedex, error)
}
