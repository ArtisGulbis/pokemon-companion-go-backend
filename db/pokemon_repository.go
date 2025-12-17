package db

import (
	"database/sql"
	"fmt"
	"log"
	"slices"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
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
	result, err := r.db.Exec(queries.InsertPokemon,
		p.ID,
		p.SpeciesID,
		p.Name,
		p.IsDefault,
		p.Height,
		p.Weight,
		p.BaseExperience,
		getStat(p.Stats, "hp"),
		getStat(p.Stats, "attack"),
		getStat(p.Stats, "defense"),
		getStat(p.Stats, "special_attack"),
		getStat(p.Stats, "special_defense"),
		getStat(p.Stats, "speed"),
		p.Sprites.Other.OfficialArtwork.FrontDefault,
		p.Sprites.Other.OfficialArtwork.FrontShiny,
		p.Sprites.Other.OfficialArtwork.FrontDefault,
	)
	if err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	rows, _ := result.RowsAffected()
	fmt.Printf("DEBUG: Inserted %d rows for pokemon ID %d\n", rows, p.ID)

	return nil
}

func (r *PokemonRepository) GetPokemonByID(id int) (*dto.Pokemon, error) {
	var pokemon dto.Pokemon

	var height, weight, baseExp, hp, attack, defense, spatk, spdef, speed sql.NullInt64
	var isDefault bool

	queryStr := `SELECT id, species_id, name, is_default, height, weight, base_experience, hp, attack, defense, special_attack, special_defense, speed, sprite_front_default, sprite_front_shiny, sprite_artwork FROM pokemon WHERE id = ?`

	err := r.db.QueryRow(queryStr, id).Scan(
		&pokemon.ID,
		&pokemon.SpeciesID,
		&pokemon.Name,
		&isDefault,
		&height,
		&weight,
		&baseExp,
		&hp,
		&attack,
		&defense,
		&spatk,
		&spdef,
		&speed,
		&pokemon.SpriteFrontDefault,
		&pokemon.SpriteFrontShiny,
		&pokemon.SpriteArtwork,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pokemon %d not found", id)
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	// Convert from database types to Go types
	pokemon.IsDefault = isDefault
	pokemon.Height = int(height.Int64)
	pokemon.Weight = int(weight.Int64)
	pokemon.BaseExperience = int(baseExp.Int64)
	pokemon.HP = int(hp.Int64)
	pokemon.Attack = int(attack.Int64)
	pokemon.Defense = int(defense.Int64)
	pokemon.SpecialAttack = int(spatk.Int64)
	pokemon.SpecialDefense = int(spdef.Int64)
	pokemon.Speed = int(speed.Int64)

	return &pokemon, nil
}

func getStat(stats []external.Stat, key string) int {
	idx := slices.IndexFunc(stats, func(c external.Stat) bool { return c.Stat.Name == key })
	if idx >= 0 {
		return stats[idx].BaseStat
	}
	return 0
}
