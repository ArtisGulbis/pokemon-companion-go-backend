package services

import (
	"testing"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockPokemonAPIClient struct {
	mock.Mock
}

func (m *MockPokemonAPIClient) FetchAll(path string) ([]external.Response, error) {
	args := m.Called(path)
	return args.Get(0).([]external.Response), args.Error(1)
}

func (m *MockPokemonAPIClient) FetchSpecies(id int) (*external.Species, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.Species), args.Error(1)
}

func (m *MockPokemonAPIClient) FetchPokemon(id int) (*external.Pokemon, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.Pokemon), args.Error(1)
}

type MockPokemonRepo struct {
	mock.Mock
}

func (m *MockPokemonRepo) InsertPokemon(p *external.Pokemon) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockPokemonRepo) InsertType(p *external.PokemonType, pokemonId int) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockPokemonRepo) InsertSpecies(p *external.Species) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockPokemonRepo) GetPokemonByID(id int) (*dto.Pokemon, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Pokemon), args.Error(1)
}

func TestSyncAllPokemon(t *testing.T) {
	t.Run("Successfully syncs all Pokemon", func(t *testing.T) {
		mockClient := new(MockPokemonAPIClient)
		mockClient.On("FetchAll", "pokemon?limit=2").Return(
			[]external.Response{
				{Name: "bulbasaur", Url: "https://pokeapi.co/api/v2/pokemon/1/"},
				{Name: "ivysaur", Url: "https://pokeapi.co/api/v2/pokemon/2/"},
			}, nil,
		)
		mockClient.On("FetchPokemon", 1).Return(
			&external.Pokemon{
				ID:     1,
				Name:   "bulbasaur",
				Height: 10,
				Weight: 100,
				Types:  []external.PokemonType{},
			}, nil,
		)

		mockClient.On("FetchPokemon", 2).Return(
			&external.Pokemon{
				ID:     2,
				Name:   "ivysaur",
				Height: 20,
				Weight: 200,
				Types:  []external.PokemonType{},
			}, nil,
		)

		mockRepo := new(MockPokemonRepo)
		mockRepo.On("InsertPokemon", mock.AnythingOfType("*external.Pokemon")).Return(nil).Twice()

		rateLimiter := time.NewTicker(1 * time.Millisecond)

		syncer := NewPokemonSyncer(mockClient, mockRepo, rateLimiter)

		err := syncer.SyncAll(2)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})
}
