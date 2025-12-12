package services

import (
	"testing"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockPokedexAPIClient struct {
	mock.Mock
}

func (m *MockPokedexAPIClient) FetchAll(path string) ([]external.Response, error) {
	args := m.Called(path)
	return args.Get(0).([]external.Response), args.Error(1)
}

func (m *MockPokedexAPIClient) FetchPokedex(id int) (*external.Pokedex, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.Pokedex), args.Error(1)
}

type MockPokedexRepo struct {
	mock.Mock
}

func (m *MockPokedexRepo) InsertPokedex(p *external.Pokedex) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockPokedexRepo) GetPokedexByID(id int) (*external.Pokedex, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.Pokedex), args.Error(1)
}

func TestSyncAllPokedexes(t *testing.T) {
	t.Run("Succesfully sync all pokedexes", func(t *testing.T) {
		mockClient := new(MockPokedexAPIClient)

		mockClient.On("FetchAll", "pokedex?limit=1").Return([]external.Response{
			{Name: "kanto", Url: "https://pokeapi.co/api/v2/pokedex/1/"},
		}, nil)

		mockClient.On("FetchPokedex", 1).Return(&external.Pokedex{
			ID:   1,
			Name: "kanto",
			Region: external.Response{
				Name: "",
				Url:  "",
			},
		}, nil)

		mockRepo := new(MockPokedexRepo)
		mockRepo.On("InsertPokedex", mock.AnythingOfType("*external.Pokedex")).Return(nil).Once()

		rateLimiter := time.NewTicker(650 * time.Millisecond)

		syncer := NewPokedexSyncer(mockClient, mockRepo, rateLimiter)

		err := syncer.SyncAll(1)
		require.NoError(t, err)
		mockClient.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestSyncPokedex(t *testing.T) {
	mockResponse := &external.Pokedex{
		ID:   1,
		Name: "johto",
		Region: external.Response{
			Name: "johto",
			Url:  "https://pokeapi.co/api/v2/region/2/",
		},
	}
	mockClient := new(MockPokedexAPIClient)
	mockClient.On("FetchPokedex", 1).Return(mockResponse, nil)

	mockRepo := new(MockPokedexRepo)
	mockRepo.On("InsertPokedex", mock.AnythingOfType("*external.Pokedex")).Return(nil).Once()

	rateLimiter := time.NewTicker(1 * time.Millisecond)

	syncer := NewPokedexSyncer(mockClient, mockRepo, rateLimiter)

	pokedex, err := syncer.SyncPokedex(mockResponse.ID)
	require.NoError(t, err)
	require.NotNil(t, pokedex)

	assert.Equal(t, mockResponse, pokedex)
}
