package services

import (
	"database/sql"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/queries"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/utils"
)

type PokemonService struct {
}

func (ps *PokemonService) InsertPokemon(p *dto.PokemonResponse) (int8, error) {
	db, _ := sql.Open("sqlite", "./pokemon.db")

	typesStmt, err := db.Prepare(queries.InsertPokemonType)
	if err != nil {
		return 1, err
	}
	defer typesStmt.Close()
	for _, p_type := range p.Types {
		type_id := utils.GetId(p_type.Type.URL)
		_, err = typesStmt.Exec(p.ID, type_id, p_type.Slot)
		if err != nil {
			return 1, err
		}
	}

	pokemonStmt, err := db.Prepare(queries.InsertPokemon)
	if err != nil {
		return 1, err
	}

	defer pokemonStmt.Close()

	_, err = pokemonStmt.Exec(p.ID,
		p.Name,
		p.Height,
		p.Weight,
	)
	if err != nil {
		return 1, err
	}
	return 0, nil
}
