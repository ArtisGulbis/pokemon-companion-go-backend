package db

import (
	"database/sql"
	"fmt"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/queries"
)

type MoveRepository struct {
	db *Database
}

func NewMoveRepository(db *Database) *MoveRepository {
	return &MoveRepository{db: db}
}

func (r *MoveRepository) InsertMove(move *external.Move) error {
	stmt, err := r.db.Prepare(queries.InsertMove)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var effect string
	if len(move.EffectEntries) > 0 {
		effect = move.EffectEntries[0].ShortEffect
	}
	_, err = stmt.Exec(
		&move.ID,
		&move.Name,
		&move.Type.Name,
		&move.Power,
		&move.Accuracy,
		&move.PP,
		&move.DamageClass.Name,
		effect,
		&move.Priority,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *MoveRepository) GetMoveByID(id int) (*dto.Move, error) {
	var move dto.Move

	err := r.db.QueryRow(queries.GetMoveByID, id).Scan(
		&move.ID,
		&move.Name,
		&move.Type,
		&move.Power,
		&move.Accuracy,
		&move.PP,
		&move.DamageClass,
		&move.EffectShort,
		&move.Priority,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("move %d not found", id)
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &move, nil
}

// InsertPokemonMove inserts a pokemon_move relationship
func (r *MoveRepository) InsertPokemonMove(pokemonID, moveID, versionGroupID int, learnMethod string, levelLearnedAt int) error {
	stmt, err := r.db.Prepare(queries.InsertPokemonMoves)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(pokemonID, moveID, versionGroupID, learnMethod, levelLearnedAt)
	if err != nil {
		return fmt.Errorf("failed to insert pokemon_move: %w", err)
	}

	return nil
}
