package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/db"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/model"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/services"
)

const pokemonBasePath = "https://pokeapi.co/api/v2/pokemon/"

type Router struct {
	db             *db.Database
	pokemonService *services.PokemonService
}

func New(database *db.Database) *Router {
	return &Router{
		db:             database,
		pokemonService: &services.PokemonService{},
	}
}

func (rt *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{id}", rt.fetchPokemon)

	// mux.HandleFunc("POST /pokemon", rt.createPokemon)
	// mux.HandleFunc("GET /pokemon", rt.listPokemon)
	// mux.HandleFunc("DELETE /{id}", rt.deletePokemon)

	return mux
}

func (rt *Router) fetchPokemon(w http.ResponseWriter, r *http.Request) {
	pokemonId := r.PathValue("id")

	resp, err := http.Get(pokemonBasePath + pokemonId)
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

	_, err = rt.pokemonService.InsertPokemon(&pokemon)
	if err != nil {
		writeError(w, "Failed to insert", err, http.StatusInternalServerError)
		return
	}

	// Set response headers and send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send the Pokemon data as JSON
	json.NewEncoder(w).Encode(map[string]*model.Pokemon{
		"pokemon": &pokemon,
	})
}

func writeError(w http.ResponseWriter, message string, err error, statusCode int) {
	log.Printf("Error: %v", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
