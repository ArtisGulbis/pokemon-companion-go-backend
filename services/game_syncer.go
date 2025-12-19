package services

import (
	"fmt"
	"log"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/utils"
)

type GameSyncer struct {
	versionSyncer *VersionSyncer
	pokedexSyncer *PokedexSyncer
	pokemonSyncer *PokemonSyncer
	rateLimiter   *time.Ticker
}

func NewGameSyncer(
	versionSyncer *VersionSyncer,
	pokedexSyncer *PokedexSyncer,
	pokemonSyncer *PokemonSyncer,
	rateLimiter *time.Ticker,
) *GameSyncer {
	return &GameSyncer{
		versionSyncer: versionSyncer,
		pokedexSyncer: pokedexSyncer,
		pokemonSyncer: pokemonSyncer,
		rateLimiter:   rateLimiter,
	}
}
func (g *GameSyncer) SyncAllGames(limit int) error {
	allVersions, err := g.versionSyncer.client.FetchAll(fmt.Sprintf("version?limit=%d", limit))
	if err != nil {
		return fmt.Errorf("failed to fetch versions: %w", err)
	}

	fmt.Printf("Found %d versions to sync", len(allVersions))

	for i, version := range allVersions {
		if i > 0 {
			<-g.rateLimiter.C
		}

		versionID, err := utils.ExtractIDFromURL(version.Url)
		if err != nil {
			return fmt.Errorf("failed to extract version ID from %s: %w", version.Url, err)
		}
		log.Printf("Syncing game %s (%d/%d)...", version.Name, i+1, len(allVersions))

		if err := g.SyncGame(versionID); err != nil {
			return fmt.Errorf("failed to sync game %d (%s): %w", versionID, version.Name, err)
		}
		log.Printf("✓ Completed %s", version.Name)
	}

	return nil
}

func (g *GameSyncer) SyncGame(id int) error {
	version, err := g.versionSyncer.client.FetchVersion(id)
	if err != nil {
		return err
	}

	versionGroupId, err := utils.ExtractIDFromURL(version.VersionGroup.Url)
	if err != nil {
		return err
	}
	versionGroup, err := g.versionSyncer.client.FetchVersionGroup(versionGroupId)
	if err != nil {
		return err
	}

	if err := g.versionSyncer.InsertVersionGroup(versionGroup); err != nil {
		return err
	}
	if err := g.versionSyncer.InsertVersion(version); err != nil {
		return err
	}
	for _, pdex := range versionGroup.Pokedexes {
		pokedexId, err := utils.ExtractIDFromURL(pdex.Url)
		if err != nil {
			return err
		}
		pokedex, err := g.pokedexSyncer.client.FetchPokedex(pokedexId)
		if err != nil {
			return err
		}
		if err := g.pokedexSyncer.InsertPokedex(pokedex); err != nil {
			return err
		}
		if err := g.pokedexSyncer.InsertVersionGroupPokedex(versionGroup); err != nil {
			return err
		}
		for _, pe := range pokedex.PokemonEntries {
			speciesID, err := utils.ExtractIDFromURL(pe.PokemonSpecies.Url)
			if err != nil {
				return err
			}
			_, err = g.pokemonSyncer.SyncSpecies(speciesID)
			if err != nil {
				return err
			}

			pokemon, err := g.pokemonSyncer.SyncPokemon(speciesID)
			if err != nil {
				return err
			}

			for _, t := range pokemon.Types {
				err = g.pokemonSyncer.InsertType(&t, pokemon.ID)
				if err != nil {
					return err
				}
			}

			for _, a := range pokemon.Abilities {
				err = g.pokemonSyncer.InsertAbility(&a, pokemon.ID)
				if err != nil {
					return err
				}
			}

			err = g.pokedexSyncer.InsertPokedexEntry(&external.PokedexEntry{
				PokedexID:   pokedexId,
				SpeciesID:   speciesID,
				EntryNumber: pe.EntryNumber,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
