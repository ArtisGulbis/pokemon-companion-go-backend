package db

import (
	"fmt"
	"log"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/queries"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/utils"
)

type PokedexRepository struct {
	db *Database
}

func NewPokedexRepository(db *Database) *PokedexRepository {
	return &PokedexRepository{db: db}
}

func (r *PokedexRepository) InsertPokedex(p *models.Pokedex) error {
	stmt, err := r.db.Prepare(queries.InsertPokedex)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		p.ID,
		p.IsMainSeries,
		p.Name,
	)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *PokedexRepository) InsertPokedexPokemonEntry(pokemonEntry []models.PokedexPokemonEntry, pokedexID int) error {
	stmt, err := r.db.Prepare(queries.InsertPokedexPokemonEntry)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, pe := range pokemonEntry {
		pokemonID, err := utils.ExtractIDFromURL(pe.PokemonSpecies.Url)
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(pe.PokemonSpecies.Name, pe.EntryNumber, pokemonID, pokedexID)
		if err != nil {
			return fmt.Errorf("failed to insert description: %w", err)
		}
	}

	return nil
}

func (r *PokedexRepository) InsertPokedexDescriptions(descriptions []models.PokedexDescriptions, pokedexID int) error {
	stmt, err := r.db.Prepare(queries.InsertPokemonDescriptions)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, desc := range descriptions {
		_, err := stmt.Exec(desc.Language.Name, desc.Description, pokedexID)
		if err != nil {
			return fmt.Errorf("failed to insert description: %w", err)
		}
	}

	return nil
}

func (r *PokedexRepository) GetPokedexEntriesByPokedexID(pokedexID int) ([]*models.PokedexPokemonEntry, error) {
	rows, err := r.db.Query(queries.GetPokedexPokemonEntryByID, pokedexID)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var pokedexPokemonEntries []*models.PokedexPokemonEntry

	for rows.Next() {
		var ppe = &models.PokedexPokemonEntry{}
		var pokemonID int
		err = rows.Scan(
			&ppe.PokemonSpecies.Name,
			&ppe.EntryNumber,
			&pokemonID,
			&pokedexID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		ppe.PokemonSpecies.Url = fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%d/", pokemonID)
		pokedexPokemonEntries = append(pokedexPokemonEntries, ppe)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	if pokedexPokemonEntries == nil {
		return nil, fmt.Errorf("pokedex pokemon entry %d not found", pokedexID)
	}
	return pokedexPokemonEntries, nil
}

func (r *PokedexRepository) GetPokedexByID(id int) (*models.Pokedex, error) {
	rows, err := r.db.Query(queries.GetPokedexByID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var pokedex *models.Pokedex

	for rows.Next() {
		pokedex = &models.Pokedex{}
		err = rows.Scan(
			&pokedex.ID,
			&pokedex.IsMainSeries,
			&pokedex.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterationg rows: %w", err)
	}
	if pokedex == nil {
		return nil, fmt.Errorf("pokedex %d not founed", id)
	}

	return pokedex, nil
}
