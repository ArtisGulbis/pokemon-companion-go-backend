package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/db"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/model"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/queries"
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

	mux.HandleFunc("POST /{id}", rt.syncPokemon)
	mux.HandleFunc("GET /{id}", rt.getPokemon)

	// mux.HandleFunc("POST /pokemon", rt.createPokemon)
	// mux.HandleFunc("GET /pokemon", rt.listPokemon)
	// mux.HandleFunc("DELETE /{id}", rt.deletePokemon)

	return mux
}

func (rt *Router) getPokemon(w http.ResponseWriter, r *http.Request) {
	pokemonId := r.PathValue("id")
	rows, err := rt.db.Query(queries.GetPokemonByID, pokemonId)
	if err != nil {
		writeError(w, "Pokemon not found", err, http.StatusNotFound)
		return
	}
	defer rows.Close()

	var pokemon *model.Pokemon

	for rows.Next() {
		var id int
		var name string
		var height int
		var weight int
		var typeName string
		var typeSlot int

		// Scan current row (order matches SELECT statement!)
		err := rows.Scan(&id, &name, &height, &weight, &typeName, &typeSlot)
		if err != nil {
			writeError(w, "Failed to read data", err, http.StatusInternalServerError)
			return
		}

		// First iteration: Create Pokemon struct
		if pokemon == nil {
			pokemon = &model.Pokemon{
				ID:     id,
				Name:   name,
				Height: height,
				Weight: weight,
				Types:  []model.PokemonType{},
			}
		}

		// Every iteration: Append type
		pokemon.Types = append(pokemon.Types, model.PokemonType{
			Slot: typeSlot,
			Name: typeName,
		})
	}

	if err := rows.Err(); err != nil {
		writeError(w, "Error reading data", err, http.StatusInternalServerError)
		return
	}

	if pokemon == nil {
		writeError(w, "Pokemon not found", nil, http.StatusNotFound)
		return
	}

	// Success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]*model.Pokemon{
		"pokemon": pokemon,
	})
}

func (rt *Router) syncPokemon(w http.ResponseWriter, r *http.Request) {
	pokemonId := r.PathValue("id")

	resp, err := http.Get(pokemonBasePath + pokemonId)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		writeError(w, "Pokemon not found", err, http.StatusNotFound)
		return
	}

	defer resp.Body.Close()

	var pokemon dto.PokemonResponse
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
	json.NewEncoder(w).Encode(map[string]*dto.PokemonResponse{
		"pokemon": &pokemon,
	})
}

func writeError(w http.ResponseWriter, message string, err error, statusCode int) {
	log.Printf("Error: %v", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
