package services

import (
	"testing"
	"time"

	models "github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
)

type MockPokedexAPIClient struct {
	FetchAllFunc     func(path string) ([]models.Response, error)
	FetchPokedexFunc func(id int) (*models.Pokedex, error)
}

func (m *MockPokedexAPIClient) FetchAll(path string) ([]models.Response, error) {
	if m.FetchAllFunc != nil {
		return m.FetchAllFunc(path)
	}
	return nil, nil
}

func (m *MockPokedexAPIClient) FetchPokedex(id int) (*models.Pokedex, error) {
	if m.FetchPokedexFunc != nil {
		return m.FetchPokedexFunc(id)
	}
	return nil, nil
}

type MockPokedexRepo struct {
	InsertPokedexFunc  func(p *models.Pokedex) error
	GetPokedexByIDFunc func(id int) (*models.Pokedex, error)
}

func (m *MockPokedexRepo) InsertPokedex(p *models.Pokedex) error {
	if m.InsertPokedexFunc != nil {
		return m.InsertPokedexFunc(p)
	}
	return nil
}

func (m *MockPokedexRepo) GetPokedexByID(id int) (*models.Pokedex, error) {
	if m.GetPokedexByIDFunc != nil {
		return m.GetPokedexByIDFunc(id)
	}
	return nil, nil
}

func TestSyncAllPokedexes(t *testing.T) {
	t.Run("Succesfully sync all pokedexes", func(t *testing.T) {
		var insertedPokedexes []*models.Pokedex

		mockClient := &MockPokedexAPIClient{
			FetchAllFunc: func(path string) ([]models.Response, error) {
				return []models.Response{
					{Name: "kanto", Url: "https://pokeapi.co/api/v2/pokedex/2/"},
					{Name: "hoenn", Url: "https://pokeapi.co/api/v2/pokedex/4/"},
				}, nil
			},
			FetchPokedexFunc: func(id int) (*models.Pokedex, error) {
				return &models.Pokedex{
					ID:   id,
					Name: "kanto",
					Region: models.Response{
						Name: "",
						Url:  "",
					},
				}, nil
			},
		}

		mockRepo := &MockPokedexRepo{
			InsertPokedexFunc: func(p *models.Pokedex) error {
				insertedPokedexes = append(insertedPokedexes, p)
				return nil
			},
		}

		rateLimiter := time.NewTicker(650 * time.Millisecond)

		syncer := NewPokedexSyncer(mockClient, mockRepo, rateLimiter)

		err := syncer.SyncAll(2)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(insertedPokedexes) != 2 {
			t.Errorf("Expected 2 Pokedexes to be inserted, got %d", len(insertedPokedexes))
		}
	})
}
