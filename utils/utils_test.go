package utils

import "testing"

func TestExtractIDFromURL(t *testing.T) {
	tests := []struct {
		name       string
		URL        string
		expectedID int
	}{
		{
			name:       "Success",
			URL:        "https://pokeapi.co/api/v2/pokedex/9",
			expectedID: 9,
		},
		{
			name:       "Success with forwardslash at the end",
			URL:        "https://pokeapi.co/api/v2/pokedex/9/",
			expectedID: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ExtractIDFromURL(tt.URL)
			if err != nil {
				t.Fatal(err)
			}
			if id != tt.expectedID {
				t.Fatalf("Expected id %d, but go %d", tt.expectedID, id)
			}
		})
	}
}
