package services

import (
	"testing"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockVersionAPIClient using testify/mock
type MockVersionAPIClient struct {
	mock.Mock
}

func (m *MockVersionAPIClient) FetchAll(path string) ([]external.Response, error) {
	args := m.Called(path)
	return args.Get(0).([]external.Response), args.Error(1)
}

func (m *MockVersionAPIClient) FetchVersionGroup(id int) (*external.VersionGroup, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.VersionGroup), args.Error(1)
}

func (m *MockVersionAPIClient) FetchVersion(id int) (*external.Version, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.Version), args.Error(1)
}

// MockVersionRepo using testify/mock
type MockVersionRepo struct {
	mock.Mock
}

func (m *MockVersionRepo) InsertVersion(v *external.Version) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockVersionRepo) InsertVersionGroup(v *external.VersionGroup) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockVersionRepo) GetVersionByID(id int) (*dto.Version, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Version), args.Error(1)
}

func TestSyncVersion(t *testing.T) {
	t.Run("Successfully sync a version", func(t *testing.T) {
		mockClient := new(MockVersionAPIClient)
		mockRepo := new(MockVersionRepo)

		mockResponse := &external.Version{
			ID:   1,
			Name: "red",
			Names: []external.Response{
				{
					Name: "red",
					Url:  "https://pokeapi.co/api/v2/language/1/",
				},
			},
			VersionGroup: external.Response{
				Name: "gold-silver",
				Url:  "https://pokeapi.co/api/v2/version-group/3/",
			},
		}

		mockClient.On("FetchVersion", 1).Return(mockResponse, nil)

		mockRepo.On("InsertVersion", mock.AnythingOfType("*external.Version")).Return(nil).Once()

		rateLimiter := time.NewTicker(1 * time.Millisecond)

		// Create syncer with mocks
		syncer := NewVersionSyncer(mockClient, mockRepo, rateLimiter)

		// Act
		response, err := syncer.SyncVersion(1)
		if err != nil {
			t.Fatalf("Failed to sync version %v", err)
		}
		require.NoError(t, err)
		require.NotNil(t, response)

		mockClient.AssertExpectations(t)
		mockRepo.AssertExpectations(t)

		assert.Equal(t, mockResponse, response)
	})
}
func TestSyncVersionGroup(t *testing.T) {
	t.Run("Successfully sync a version group", func(t *testing.T) {
		mockClient := new(MockVersionAPIClient)
		mockRepo := new(MockVersionRepo)

		mockResponse := &external.VersionGroup{
			ID:   1,
			Name: "red",
			Generation: external.Response{
				Name: "generation-i",
				Url:  "https://pokeapi.co/api/v2/version-group/1/",
			},
			Pokedexes: []external.Response{
				{
					Name: "original-johto",
					Url:  "https://pokeapi.co/api/v2/pokedex/3/",
				},
			},
		}

		mockClient.On("FetchVersionGroup", 1).Return(mockResponse, nil)

		mockRepo.On("InsertVersionGroup", mock.AnythingOfType("*external.VersionGroup")).Return(nil).Once()

		rateLimiter := time.NewTicker(1 * time.Millisecond)

		// Create syncer with mocks
		syncer := NewVersionSyncer(mockClient, mockRepo, rateLimiter)

		// Act
		response, err := syncer.SyncVersionGroup(1)
		if err != nil {
			t.Fatalf("Failed to sync version group %v", err)
		}
		require.NoError(t, err)
		require.NotNil(t, response)

		mockClient.AssertExpectations(t)
		mockRepo.AssertExpectations(t)

		assert.Equal(t, mockResponse, response)
	})
}

func TestSyncAllVersions(t *testing.T) {
	t.Run("Successfully syncs all versions", func(t *testing.T) {
		// Create mocks
		mockClient := new(MockVersionAPIClient)
		mockRepo := new(MockVersionRepo)

		// Set up expectations - what we expect to be called
		mockClient.On("FetchAll", "version?limit=2").Return([]external.Response{
			{Name: "red", Url: "https://pokeapi.co/api/v2/version/1/"},
			{Name: "blue", Url: "https://pokeapi.co/api/v2/version/2/"},
		}, nil)

		// Expect FetchVersion to be called twice (for ids 1 and 2)
		mockClient.On("FetchVersion", 1).Return(&external.Version{
			ID:   1,
			Name: "red",
			VersionGroup: external.Response{
				Name: "red-blue",
				Url:  "https://pokeapi.co/api/v2/version-group/1/",
			},
		}, nil)

		mockClient.On("FetchVersion", 2).Return(&external.Version{
			ID:   2,
			Name: "blue",
			VersionGroup: external.Response{
				Name: "red-blue",
				Url:  "https://pokeapi.co/api/v2/version-group/1/",
			},
		}, nil)

		// Expect InsertVersion to be called twice
		mockRepo.On("InsertVersion", mock.AnythingOfType("*external.Version")).Return(nil).Twice()

		rateLimiter := time.NewTicker(1 * time.Millisecond)

		// Create syncer with mocks
		syncer := NewVersionSyncer(mockClient, mockRepo, rateLimiter)

		// Act
		err := syncer.SyncAll(2)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify all expectations were met
		mockClient.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
