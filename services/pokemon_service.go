package services

import (
	"database/sql"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/model"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/queries"
)

type PokemonService struct {
}

func (ps *PokemonService) InsertPokemon(p *model.Pokemon) (int8, error) {
	db, _ := sql.Open("sqlite", "./pokemon.db")
	stmt, err := db.Prepare(queries.InsertPokemon)
	if err != nil {
		return 1, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(p.ID,
		p.Name,
		p.Height,
		p.Weight,
	)
	if err != nil {
		return 1, err
	}
	return 0, nil
}
