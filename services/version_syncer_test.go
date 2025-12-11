package services

import (
	"testing"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	models "github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/mock"
)

// MockVersionAPIClient using testify/mock
type MockVersionAPIClient struct {
	mock.Mock
}

func (m *MockVersionAPIClient) FetchAll(path string) ([]models.Response, error) {
	args := m.Called(path)
	return args.Get(0).([]models.Response), args.Error(1)
}

func (m *MockVersionAPIClient) FetchVersion(id int) (*models.Version, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Version), args.Error(1)
}

// MockVersionRepo using testify/mock
type MockVersionRepo struct {
	mock.Mock
}

func (m *MockVersionRepo) InsertVersion(v *external.Version) error {
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

func TestSyncAllVersions(t *testing.T) {
	t.Run("Successfully syncs all versions", func(t *testing.T) {
		// Create mocks
		mockClient := new(MockVersionAPIClient)
		mockRepo := new(MockVersionRepo)

		// Set up expectations - what we expect to be called
		mockClient.On("FetchAll", "version?limit=2").Return([]models.Response{
			{Name: "red", Url: "https://pokeapi.co/api/v2/version/1/"},
			{Name: "blue", Url: "https://pokeapi.co/api/v2/version/2/"},
		}, nil)

		// Expect FetchVersion to be called twice (for ids 1 and 2)
		mockClient.On("FetchVersion", 1).Return(&models.Version{
			ID:   1,
			Name: "red",
			VersionGroup: models.Response{
				Name: "red-blue",
				Url:  "https://pokeapi.co/api/v2/version-group/1/",
			},
		}, nil)

		mockClient.On("FetchVersion", 2).Return(&models.Version{
			ID:   2,
			Name: "blue",
			VersionGroup: models.Response{
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
