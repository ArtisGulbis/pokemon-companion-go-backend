package services

import (
	"testing"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/mock"
)

type MockMoveAPIClient struct {
	mock.Mock
}

func (m *MockMoveAPIClient) FetchAll(path string) ([]external.Response, error) {
	args := m.Called(path)
	return args.Get(0).([]external.Response), args.Error(1)
}

func (m *MockMoveAPIClient) FetchMove(id int) (*external.Move, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.Move), args.Error(1)
}

type MockMoveRepo struct {
	mock.Mock
}

func (mr *MockMoveRepo) InsertMove(m *external.Move) error {
	args := mr.Called(m)
	return args.Error(0)
}

func (mr *MockMoveRepo) InsertPokemonMove(pokemonID int, moveID int, versionGroupID int, learnMethod string, levelLearnedAt int) error {
	args := mr.Called(pokemonID, moveID, versionGroupID, learnMethod, levelLearnedAt)
	return args.Error(0)
}

func (mr *MockMoveRepo) GetMoveByID(id int) (*dto.Move, error) {
	args := mr.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Move), args.Error(1)
}

func TestSyncMove(t *testing.T) {
	t.Run("Sync move success", func(t *testing.T) {
		mockClient := new(MockMoveAPIClient)
		mockRepo := new(MockMoveRepo)

		mockResponse := &external.Move{
			ID:   1,
			Name: "tackle",
			Type: external.Response{
				Name: "normal",
				Url:  "",
			},
			Power:    100,
			Accuracy: 100,
			PP:       10,
			DamageClass: external.Response{
				Name: "physical",
				Url:  "",
			},
			EffectEntries: []external.EffectEntry{
				{
					Effect: "",
					Language: external.Response{
						Name: "",
						Url:  "",
					},
					ShortEffect: "hit the enemy",
				},
			},
			Priority: 0,
		}

		mockClient.On("FetchMove", 1).Return(mockResponse, nil)
		mockRepo.On("InsertMove", mock.AnythingOfType("*external.Move")).Return(nil).Once()

		rateLimiter := time.NewTicker(1 * time.Millisecond)

		syncer := NewMoveSyncer(mockClient, mockRepo, rateLimiter)
		err := syncer.SyncMove(1)
		if err != nil {
			t.Fatal(err)
		}
		// require.NotNil(t, response)
		mockClient.AssertExpectations(t)
		mockRepo.AssertExpectations(t)

		// assert.Equal(t, mockResponse, response)
	})
}
