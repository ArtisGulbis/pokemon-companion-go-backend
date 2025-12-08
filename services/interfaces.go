package services

import models "github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"

type PokemonAPIClient interface {
	FetchAll(path string) ([]models.Response, error)
	FetchPokemon(id int) (*models.Pokemon, error)
}

type PokedexAPIClient interface {
	FetchAll(path string) ([]models.Response, error)
	FetchPokedex(id int) (*models.Pokedex, error)
}

type PokemonRepo interface {
	InsertPokemon(p *models.Pokemon) error
	GetPokemonByID(id int) (*models.Pokemon, error)
}

type PokedexRepo interface {
	InsertPokedex(p *models.Pokedex) error
	GetPokedexByID(id int) (*models.Pokedex, error)
}
