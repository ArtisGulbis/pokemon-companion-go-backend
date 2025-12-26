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
	moveSyncer    *MoveSyncer
	rateLimiter   *time.Ticker
}

func NewGameSyncer(
	versionSyncer *VersionSyncer,
	pokedexSyncer *PokedexSyncer,
	pokemonSyncer *PokemonSyncer,
	moveSyncer *MoveSyncer,
	rateLimiter *time.Ticker,
) *GameSyncer {
	return &GameSyncer{
		versionSyncer: versionSyncer,
		pokedexSyncer: pokedexSyncer,
		pokemonSyncer: pokemonSyncer,
		moveSyncer:    moveSyncer,
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

	log.Printf("Inserting version group %d (%s)...", versionGroup.ID, versionGroup.Name)
	if err := g.versionSyncer.InsertVersionGroup(versionGroup); err != nil {
		return fmt.Errorf("failed to insert version group %d (%s): %w", versionGroup.ID, versionGroup.Name, err)
	}

	log.Printf("Inserting version %d (%s) with version_group_id %d...", version.ID, version.Name, versionGroupId)
	if err := g.versionSyncer.InsertVersion(version); err != nil {
		return fmt.Errorf("failed to insert version %d (%s): %w", version.ID, version.Name, err)
	}
	log.Printf("Version %d inserted successfully", version.ID)

	log.Printf("Processing %d pokedexes for version group %d...", len(versionGroup.Pokedexes), versionGroup.ID)
	for i, pdex := range versionGroup.Pokedexes {
		pokedexId, err := utils.ExtractIDFromURL(pdex.Url)
		if err != nil {
			return fmt.Errorf("failed to extract pokedex ID: %w", err)
		}
		log.Printf("  [%d/%d] Fetching pokedex %d...", i+1, len(versionGroup.Pokedexes), pokedexId)
		pokedex, err := g.pokedexSyncer.client.FetchPokedex(pokedexId)
		if err != nil {
			return fmt.Errorf("failed to fetch pokedex %d: %w", pokedexId, err)
		}
		log.Printf("  [%d/%d] Inserting pokedex %d (%s)...", i+1, len(versionGroup.Pokedexes), pokedex.ID, pokedex.Name)
		if err := g.pokedexSyncer.InsertPokedex(pokedex); err != nil {
			return fmt.Errorf("failed to insert pokedex %d: %w", pokedex.ID, err)
		}
	}

	// Insert version_group_pokedex relationships AFTER all pokedexes are inserted
	log.Printf("Inserting version_group_pokedex relationships...")
	if err := g.pokedexSyncer.InsertVersionGroupPokedex(versionGroup); err != nil {
		return fmt.Errorf("failed to insert version_group_pokedex for vg %d: %w", versionGroup.ID, err)
	}

	// Process Pokemon entries for each pokedex
	for _, pdex := range versionGroup.Pokedexes {
		pokedexId, err := utils.ExtractIDFromURL(pdex.Url)
		if err != nil {
			return fmt.Errorf("failed to extract pokedex ID: %w", err)
		}
		pokedex, err := g.pokedexSyncer.client.FetchPokedex(pokedexId)
		if err != nil {
			return fmt.Errorf("failed to fetch pokedex %d: %w", pokedexId, err)
		}
		log.Printf("Processing %d Pokemon entries for pokedex %d (%s)...", len(pokedex.PokemonEntries), pokedex.ID, pokedex.Name)
		for _, pe := range pokedex.PokemonEntries {
			speciesID, err := utils.ExtractIDFromURL(pe.PokemonSpecies.Url)
			if err != nil {
				return fmt.Errorf("failed to extract species ID: %w", err)
			}

			log.Printf("  Syncing species %d...", speciesID)
			_, err = g.pokemonSyncer.SyncSpecies(speciesID)
			if err != nil {
				return fmt.Errorf("failed to sync species %d: %w", speciesID, err)
			}

			log.Printf("  Syncing pokemon %d...", speciesID)
			pokemon, err := g.pokemonSyncer.SyncPokemon(speciesID)
			if err != nil {
				return fmt.Errorf("failed to sync pokemon %d: %w", speciesID, err)
			}

			log.Printf("  Inserting %d types for pokemon %d (%s)...", len(pokemon.Types), pokemon.ID, pokemon.Name)
			for _, t := range pokemon.Types {
				err = g.pokemonSyncer.InsertType(&t, pokemon.ID)
				if err != nil {
					return fmt.Errorf("failed to insert type %s for pokemon %d: %w", t.Type.Name, pokemon.ID, err)
				}
			}

			log.Printf("  Syncing %d moves for pokemon %d (%s)...", len(pokemon.Moves), pokemon.ID, pokemon.Name)
			for _, m := range pokemon.Moves {
				moveId, err := utils.ExtractIDFromURL(m.Move.Url)
				if err != nil {
					return fmt.Errorf("failed to extract move ID: %w", err)
				}
				err = g.moveSyncer.SyncMove(moveId)
				if err != nil {
					return fmt.Errorf("failed to sync move %d for pokemon %d: %w", moveId, pokemon.ID, err)
				}
			}

			log.Printf("  Inserting %d abilities for pokemon %d (%s)...", len(pokemon.Abilities), pokemon.ID, pokemon.Name)
			for _, a := range pokemon.Abilities {
				err = g.pokemonSyncer.InsertAbility(&a, pokemon.ID)
				if err != nil {
					return fmt.Errorf("failed to insert ability %s for pokemon %d: %w", a.Ability.Name, pokemon.ID, err)
				}
			}

			log.Printf("  Inserting pokedex entry for species %d...", speciesID)
			err = g.pokedexSyncer.InsertPokedexEntry(&external.PokedexEntry{
				PokedexID:   pokedexId,
				SpeciesID:   speciesID,
				EntryNumber: pe.EntryNumber,
			})
			if err != nil {
				return fmt.Errorf("failed to insert pokedex entry for species %d: %w", speciesID, err)
			}
		}
	}

	return nil
}
