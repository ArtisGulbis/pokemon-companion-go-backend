package services

import (
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models"
)

// MockPokemonAPIClient is a mock implementation of PokemonAPIClient
type MockPokemonAPIClient struct {
	FetchAllFunc     func(path string) ([]models.Response, error)
	FetchPokemonFunc func(id int) (*models.Pokemon, error)
}

func (m *MockPokemonAPIClient) FetchAll(path string) ([]models.Response, error) {
	if m.FetchAllFunc != nil {
		return m.FetchAllFunc(path)
	}
	return nil, nil
}

func (m *MockPokemonAPIClient) FetchPokemon(id int) (*models.Pokemon, error) {
	if m.FetchPokemonFunc != nil {
		return m.FetchPokemonFunc(id)
	}
	return nil, nil
}

// MockPokemonRepo is a mock implementation of PokemonRepo
type MockPokemonRepo struct {
	InsertPokemonFunc  func(p *models.Pokemon) error
	GetPokemonByIDFunc func(id int) (*models.Pokemon, error)
}

func (m *MockPokemonRepo) InsertPokemon(p *models.Pokemon) error {
	if m.InsertPokemonFunc != nil {
		return m.InsertPokemonFunc(p)
	}
	return nil
}

func (m *MockPokemonRepo) GetPokemonByID(id int) (*models.Pokemon, error) {
	if m.GetPokemonByIDFunc != nil {
		return m.GetPokemonByIDFunc(id)
	}
	return nil, nil
}

func TestSyncAll(t *testing.T) {
	t.Run("Successfully syncs all Pokemon", func(t *testing.T) {
		// Track what was called
		var insertedPokemon []*models.Pokemon

		// Create mock client
		mockClient := &MockPokemonAPIClient{
			FetchAllFunc: func(path string) ([]models.Response, error) {
				return []models.Response{
					{Name: "bulbasaur", Url: "https://pokeapi.co/api/v2/pokemon/1/"},
					{Name: "ivysaur", Url: "https://pokeapi.co/api/v2/pokemon/2/"},
				}, nil
			},
			FetchPokemonFunc: func(id int) (*models.Pokemon, error) {
				// Return mock Pokemon based on ID
				return &models.Pokemon{
					ID:     id,
					Name:   "test-pokemon",
					Height: 10,
					Weight: 100,
					Types:  []models.PokemonType{},
				}, nil
			},
		}

		// Create mock repo
		mockRepo := &MockPokemonRepo{
			InsertPokemonFunc: func(p *models.Pokemon) error {
				// Track what was inserted
				insertedPokemon = append(insertedPokemon, p)
				return nil
			},
		}

		// Create syncer with mocks
		syncer := NewPokemonSyncer(mockClient, mockRepo)

		// Act
		err := syncer.SyncAll(2)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify InsertPokemon was called twice
		if len(insertedPokemon) != 2 {
			t.Errorf("Expected 2 Pokemon to be inserted, got %d", len(insertedPokemon))
		}
	})
}
