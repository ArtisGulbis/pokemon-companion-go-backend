package services

import "github.com/ArtisGulbis/pokemon-companion-go-backend/models"

// PokemonAPIClient defines the interface for fetching Pokemon data
type PokemonAPIClient interface {
	FetchAll(path string) ([]models.Response, error)
	FetchPokemon(id int) (*models.Pokemon, error)
}

// PokemonRepo defines the interface for Pokemon database operations
type PokemonRepo interface {
	InsertPokemon(p *models.Pokemon) error
	GetPokemonByID(id int) (*models.Pokemon, error)
}
