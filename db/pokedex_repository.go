package db

import (
	"fmt"
	"log"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/queries"
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
