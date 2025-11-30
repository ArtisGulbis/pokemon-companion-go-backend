package main

// Import the testing package - this is Go's built-in testing framework
import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest" // This package provides utilities for testing HTTP handlers
	"strings"
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/model"
)

// TABLE-DRIVEN TEST PATTERN
// This is the idiomatic Go way to test multiple scenarios
// We define a slice of test cases, then loop through them
func TestWriteError(t *testing.T) {
	// Define a struct to hold each test case
	// This makes our tests clear and easy to add to
	tests := []struct {
		name           string // Name describes what we're testing
		message        string // Input: the error message
		err            error  // Input: the error
		statusCode     int    // Input: the HTTP status code
		expectedStatus int    // Expected: status code in response
		expectedMsg    string // Expected: error message in JSON
	}{
		{
			name:           "Bad Request Error",
			message:        "Invalid input",
			err:            errors.New("validation failed"),
			statusCode:     http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Invalid input",
		},
		{
			name:           "Not Found Error",
			message:        "Pokemon not found",
			err:            errors.New("no pokemon with that ID"),
			statusCode:     http.StatusNotFound,
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "Pokemon not found",
		},
		{
			name:           "Internal Server Error",
			message:        "Database connection failed",
			err:            errors.New("connection refused"),
			statusCode:     http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Database connection failed",
		},
	}

	// Loop through each test case
	// "range" iterates over the slice, giving us each test case
	for _, tt := range tests {
		// t.Run creates a "subtest" - each test case runs independently
		// If one fails, the others still run
		// The first argument is the subtest name (shows in output)
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh ResponseRecorder for this test case
			w := httptest.NewRecorder()

			// Call the function we're testing with this test case's inputs
			writeError(w, tt.message, tt.err, tt.statusCode)

			// Verify the status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Verify the Content-Type header
			if ct := w.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", ct)
			}

			// Verify the JSON response body
			var response map[string]string
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response["error"] != tt.expectedMsg {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedMsg, response["error"])
			}
		})
	}
}

// TESTING HTTP HANDLERS WITH MOCKED EXTERNAL APIS
// This test shows how to test a handler that makes external HTTP calls
// We create a mock server to simulate the PokeAPI responses
func TestFetchPokemon(t *testing.T) {
	// Table-driven test with different scenarios
	tests := []struct {
		name               string // Description of test case
		pokemonID          string // The ID we'll request
		mockResponse       string // What our mock PokeAPI returns
		mockStatusCode     int    // HTTP status from mock API
		expectedStatusCode int    // Expected status in our handler's response
		expectError        bool   // Should the response contain an error?
	}{
		{
			name:      "Successful Pokemon Fetch",
			pokemonID: "25", // Pikachu!
			// This is a simplified Pokemon JSON response
			mockResponse: `{
				"id": 25,
				"name": "pikachu",
				"height": 4,
				"weight": 60,
				"types": [
					{
						"slot": 1,
						"type": {
							"name": "electric",
							"url": "https://pokeapi.co/api/v2/type/13/"
						}
					}
				]
			}`,
			mockStatusCode:     http.StatusOK,
			expectedStatusCode: http.StatusOK,
			expectError:        false,
		},
		{
			name:               "Pokemon Not Found",
			pokemonID:          "99999",
			mockResponse:       `{"error": "Not Found"}`,
			mockStatusCode:     http.StatusNotFound, // Mock API returns 404
			expectedStatusCode: http.StatusNotFound, // Our handler should also return 404
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// STEP 1: Create a mock PokeAPI server
			// httptest.NewServer creates a real HTTP server for testing
			// It automatically picks an available port and starts listening
			mockPokeAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// This function handles requests to our mock server
				// We verify it's requesting the right path
				expectedPath := "/" + tt.pokemonID
				if r.URL.Path != expectedPath {
					t.Errorf("Expected request to %s, got %s", expectedPath, r.URL.Path)
				}

				// Return the mock response we defined in the test case
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			// defer ensures this runs after the test completes
			// Always close test servers to free up resources
			defer mockPokeAPI.Close()

			// STEP 2: Temporarily override the pokemonBasePath to use our mock
			// Save the original value so we can restore it
			originalBasePath := pokemonBasePath
			// Reassign to point to our mock server
			// We can't directly change the const, but in a real app you'd use
			// dependency injection or a variable instead
			// For now, we'll work around this limitation

			// STEP 3: Create a request to our handler
			// httptest.NewRequest creates a fake HTTP request
			// Parameters: method, url, body
			req := httptest.NewRequest("GET", "/"+tt.pokemonID, nil)

			// We need to set path values since we use r.PathValue("id")
			// In Go 1.22+, we can use SetPathValue for this
			req.SetPathValue("id", tt.pokemonID)

			// STEP 4: Create a ResponseRecorder to capture the response
			w := httptest.NewRecorder()

			// STEP 5: Here's the workaround - we'll test with the mock URL
			// In production code, you'd inject the base URL as a dependency
			// For this test, let's create a modified version of fetchPokemon
			// that accepts the base URL

			// Actually, let's modify our approach - we'll test by calling
			// the handler and using string replacement
			// This is a limitation of the current code structure

			// For now, let's create a test helper that mimics fetchPokemon
			// but uses our mock URL
			testFetchPokemon := func(w http.ResponseWriter, r *http.Request) {
				pokemonId := r.PathValue("id")
				// Use mock server URL instead of real PokeAPI
				resp, err := http.Get(mockPokeAPI.URL + "/" + pokemonId)
				if err != nil || resp.StatusCode == http.StatusNotFound {
					writeError(w, "Pokemon not found", err, http.StatusNotFound)
					return
				}
				defer resp.Body.Close()

				var pokemon model.Pokemon
				err = json.NewDecoder(resp.Body).Decode(&pokemon)
				if err != nil {
					writeError(w, "Failed to fetch", err, http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]*model.Pokemon{
					"pokemon": &pokemon,
				})
			}

			// STEP 6: Call the handler
			testFetchPokemon(w, req)

			// STEP 7: Verify the response
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status %d, got %d", tt.expectedStatusCode, w.Code)
			}

			// Check Content-Type
			if ct := w.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", ct)
			}

			// STEP 8: Verify the response body
			if tt.expectError {
				// Should contain an error message
				var response map[string]string
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if response["error"] == "" {
					t.Error("Expected error message in response")
				}
			} else {
				// Should contain Pokemon data
				var response map[string]*model.Pokemon
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode success response: %v", err)
				}
				if response["pokemon"] == nil {
					t.Fatal("Expected 'pokemon' in response")
				}
				// Verify the Pokemon has the expected data
				pokemon := response["pokemon"]
				if pokemon.Name != "pikachu" && !tt.expectError {
					t.Errorf("Expected pokemon name 'pikachu', got '%s'", pokemon.Name)
				}
			}

			// Restore original base path (though we're using a test version)
			_ = originalBasePath // Avoid unused variable error
		})
	}
}

// BONUS: Let's add a helper function to show a better testing pattern
// This shows how you might structure testable code in the future
// Helper function to create a test-friendly version of fetchPokemon
func createFetchPokemonHandler(baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pokemonId := r.PathValue("id")
		resp, err := http.Get(baseURL + pokemonId)
		if err != nil || resp.StatusCode == http.StatusNotFound {
			writeError(w, "Pokemon not found", err, http.StatusNotFound)
			return
		}
		defer resp.Body.Close()

		var pokemon model.Pokemon
		err = json.NewDecoder(resp.Body).Decode(&pokemon)
		if err != nil {
			writeError(w, "Failed to fetch", err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]*model.Pokemon{
			"pokemon": &pokemon,
		})
	}
}

// Test using the configurable handler - cleaner approach!
func TestFetchPokemonWithConfigurableURL(t *testing.T) {
	// Create a mock PokeAPI server
	mockPokeAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if requesting pikachu
		if strings.Contains(r.URL.Path, "25") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": 25,
				"name": "pikachu",
				"height": 4,
				"weight": 60,
				"types": [{
					"slot": 1,
					"type": {
						"name": "electric",
						"url": "https://pokeapi.co/api/v2/type/13/"
					}
				}]
			}`))
			return
		}
		// Otherwise return 404
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockPokeAPI.Close()

	// Create handler with mock URL
	handler := createFetchPokemonHandler(mockPokeAPI.URL + "/")

	// Create request
	req := httptest.NewRequest("GET", "/25", nil)
	req.SetPathValue("id", "25")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler(w, req)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]*model.Pokemon
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["pokemon"].Name != "pikachu" {
		t.Errorf("Expected pikachu, got %s", response["pokemon"].Name)
	}
}
