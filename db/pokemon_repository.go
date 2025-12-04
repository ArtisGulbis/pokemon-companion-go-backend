package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/queries"
)

type PokemonRepository struct {
	db *Database
}

func NewPokemonRepository(db *Database) *PokemonRepository {
	return &PokemonRepository{db: db}
}

func (r *PokemonRepository) InsertPokemon(p *models.Pokemon) error {
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
	)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (r *PokemonRepository) GetPokemonByID(id int) (*models.Pokemon, error) {
	rows, err := r.db.Query(queries.GetPokemonByID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var pokemon *models.Pokemon

	for rows.Next() {
		var typeName sql.NullString
		var typeSlot sql.NullInt64

		if pokemon == nil {
			pokemon = &models.Pokemon{Types: []models.PokemonType{}}
			err = rows.Scan(
				&pokemon.ID,
				&pokemon.Name,
				&pokemon.Height,
				&pokemon.Weight,
				&typeName,
				&typeSlot,
			)
		} else {
			var tempID int
			var tempName string
			var tempHeight, tempWeight int
			err = rows.Scan(&tempID, &tempName, &tempHeight, &tempWeight, &typeName, &typeSlot)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if typeName.Valid && typeSlot.Valid {
			pokemon.Types = append(pokemon.Types, models.PokemonType{
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
