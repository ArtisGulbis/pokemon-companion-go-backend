package pokeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models"
)

func TestFetchPokemon(t *testing.T) {
	tests := []struct {
		name           string
		pokemonID      int
		mockStatus     int
		mockResponse   string
		expectErr      bool
		expectedID     int
		expectedName   string
		expectedHeight int
		expectedWeight int
	}{
		{
			name:           "Valid Pokemon",
			pokemonID:      25,
			mockStatus:     200,
			mockResponse:   `{"id": 25, "name": "pikachu", "height": 4, "weight": 60}`,
			expectErr:      false,
			expectedID:     25,
			expectedName:   "pikachu",
			expectedHeight: 4,
			expectedWeight: 60,
		},
		{
			name:         "Not Found",
			pokemonID:    99999,
			mockStatus:   404,
			mockResponse: ``,
			expectErr:    true,
			expectedID:   0,
			expectedName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/api/v2/pokemon/%d", tt.pokemonID) {
					t.Errorf("Expected path /api/v2/pokemon/%d, got %s", tt.pokemonID, r.URL.Path)
				}
				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()
			client := NewClient(mockServer.URL)
			pokemon, err := client.FetchPokemon(tt.pokemonID)
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return // Don't check pokemon fields if we expected error
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if pokemon.ID != tt.expectedID {
				t.Errorf("Expected ID %d , got %d", tt.expectedID, pokemon.ID)
			}
			if pokemon.Name != tt.expectedName {
				t.Errorf("Expected Name %s, got %s", tt.expectedName, pokemon.Name)
			}
			if pokemon.Weight != tt.expectedWeight {
				t.Errorf("Expected Weight %d, got %d", tt.expectedWeight, pokemon.Weight)
			}
			if pokemon.Height != tt.expectedHeight {
				t.Errorf("Expected Heigth %d, got %d", tt.expectedHeight, pokemon.Height)
			}
		})
	}
}

func TestFetchAllPokemon(t *testing.T) {
	tests := []struct {
		name             string
		mockResponse     string
		mockStatus       int
		expectedResponse []models.Response
	}{
		{
			name:       "Not Found",
			mockStatus: 200,
			mockResponse: `{
				"results": [
					{
					"name": "bulbasaur",
					"url": "https://pokeapi.co/api/v2/pokemon/1/"
					},
					{
					"name": "ivysaur",
					"url": "https://pokeapi.co/api/v2/pokemon/2/"
					}
				]
			}`,
			expectedResponse: []models.Response{
				{Name: "bulbasaur", Url: "https://pokeapi.co/api/v2/pokemon/1/"},
				{Name: "ivysaur", Url: "https://pokeapi.co/api/v2/pokemon/2/"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v2/pokemon" {
					t.Errorf("Expected path /api/v2/pokemon?limit=1500, got %s", r.URL.Path)
				}
				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()
			client := NewClient(mockServer.URL)
			allPokemon, err := client.FetchAllPokemon()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			for idx, p := range allPokemon {
				current := tt.expectedResponse[idx]
				if p.Name != current.Name {
					t.Errorf("Expected Name %s, got %s", current.Name, p.Name)
				}
				if p.Url != current.Url {
					t.Errorf("Expected Url %s, got %s", current.Url, p.Name)
				}
			}

		})
	}
}
