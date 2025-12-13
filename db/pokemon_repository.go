package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/queries"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/utils"
)

type PokemonRepository struct {
	db *Database
}

func NewPokemonRepository(db *Database) *PokemonRepository {
	return &PokemonRepository{db: db}
}

func (r *PokemonRepository) InsertSpecies(s *external.Species) error {
	stmt, err := r.db.Prepare(queries.InsertSpecies)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	evolutionChainID, err := utils.ExtractIDFromURL(s.EvolutionChain.URL)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		s.ID,
		s.Name,
		evolutionChainID,
		s.GenderRate,
		s.CaptureRate,
		s.BaseHappiness,
		s.IsBaby,
		s.IsLegendary,
		s.IsMythical,
		s.GrowthRate.Name,
		s.Generation.Name,
	)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (r *PokemonRepository) InsertPokemon(p *external.Pokemon) error {
	stmt, err := r.db.Prepare(queries.InsertPokemon)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		p.ID,
		p.Name,
		p.Height,
		p.Weight,
		p.BaseExperience,
		p.SpeciesID,
		p.IsDefault,
	)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (r *PokemonRepository) GetPokemonByID(id int) (*external.Pokemon, error) {
	rows, err := r.db.Query(queries.GetPokemonByID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var pokemon *external.Pokemon

	for rows.Next() {
		var typeName sql.NullString
		var typeSlot sql.NullInt64
		var baseExperience sql.NullInt64

		if pokemon == nil {
			pokemon = &external.Pokemon{Types: []external.PokemonType{}}
			err = rows.Scan(
				&pokemon.ID,
				&pokemon.Name,
				&pokemon.Height,
				&pokemon.Weight,
				&baseExperience,
				&typeName,
				&typeSlot,
			)
			if baseExperience.Valid {
				pokemon.BaseExperience = int(baseExperience.Int64)
			}
		} else {
			var tempID int
			var tempName string
			var tempHeight, tempWeight int
			var tempBaseExperience sql.NullInt64
			err = rows.Scan(&tempID, &tempName, &tempHeight, &tempWeight, &tempBaseExperience, &typeName, &typeSlot)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if typeName.Valid && typeSlot.Valid {
			pokemon.Types = append(pokemon.Types, external.PokemonType{
				Name: typeName.String,
				Slot: int(typeSlot.Int64),
			})
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if pokemon == nil {
		return nil, fmt.Errorf("pokemon %d not found", id)
	}

	return pokemon, nil
}
