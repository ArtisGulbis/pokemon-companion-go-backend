package services

import (
	"fmt"
	"log"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models"
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

	if version.Name == "green-japan" || version.Name == "red-japan" || version.Name == "blue-japan" {
		return nil
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

	// Check if this is a special game version (Colosseum, XD, etc.) that doesn't have traditional Pokedexes
	if specialPokemonIDs := models.GetSpecialGamePokemon(version.Name); specialPokemonIDs != nil {
		log.Printf("Special game detected: %s - Processing %d Pokemon...", version.Name, len(specialPokemonIDs))
		return g.syncSpecialGamePokemon(specialPokemonIDs, version.Name, versionGroupId)
	}

	log.Printf("Processing %d pokedexes for version group %d...", len(versionGroup.Pokedexes), versionGroup.ID)

	// Cache fetched pokedexes to avoid duplicate fetches
	pokedexCache := make(map[int]*external.Pokedex)

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

		// Store in cache for later use
		pokedexCache[pokedexId] = pokedex

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

	// Process Pokemon entries for each pokedex using cached data
	for _, pdex := range versionGroup.Pokedexes {
		pokedexId, err := utils.ExtractIDFromURL(pdex.Url)
		if err != nil {
			return fmt.Errorf("failed to extract pokedex ID: %w", err)
		}

		// Use cached pokedex instead of fetching again
		pokedex, ok := pokedexCache[pokedexId]
		if !ok {
			return fmt.Errorf("pokedex %d not found in cache", pokedexId)
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

			log.Printf("  Syncing pokemon %d with types, moves, and abilities...", speciesID)
			if err := g.syncPokemonData(speciesID, versionGroupId); err != nil {
				return fmt.Errorf("failed to sync pokemon %d: %w", speciesID, err)
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

// syncSpecialGamePokemon handles syncing Pokemon for special game versions like Colosseum and XD
// that don't have traditional Pokedexes in PokeAPI
func (g *GameSyncer) syncSpecialGamePokemon(pokemonIDs []int, versionName string, versionGroupID int) error {
	log.Printf("Syncing %d Pokemon for special game version...", len(pokemonIDs))

	// Create a virtual pokedex for this special game
	// Use a high ID that won't conflict with real pokedexes (real ones are < 100)
	virtualPokedexID := 1000 + versionGroupID
	virtualPokedex := &external.Pokedex{
		ID:   virtualPokedexID,
		Name: versionName + "-pokedex",
		Region: external.Response{
			Name: "unknown", // Special games don't have a specific region
		},
	}

	log.Printf("Creating virtual pokedex %d (%s)...", virtualPokedex.ID, virtualPokedex.Name)
	if err := g.pokedexSyncer.InsertPokedex(virtualPokedex); err != nil {
		return fmt.Errorf("failed to insert virtual pokedex: %w", err)
	}

	// Link the virtual pokedex to the version group
	versionGroup := &external.VersionGroup{
		ID: versionGroupID,
		Pokedexes: []external.Response{
			{Name: virtualPokedex.Name, Url: fmt.Sprintf("https://pokeapi.co/api/v2/pokedex/%d/", virtualPokedexID)},
		},
	}
	if err := g.pokedexSyncer.InsertVersionGroupPokedex(versionGroup); err != nil {
		return fmt.Errorf("failed to insert version_group_pokedex: %w", err)
	}

	// Sync each Pokemon and create pokedex entries
	for i, pokemonID := range pokemonIDs {
		log.Printf("  [%d/%d] Processing Pokemon %d...", i+1, len(pokemonIDs), pokemonID)

		// Sync species
		_, err := g.pokemonSyncer.SyncSpecies(pokemonID)
		if err != nil {
			return fmt.Errorf("failed to sync species %d: %w", pokemonID, err)
		}

		// Sync Pokemon and get full data
		if err := g.syncPokemonData(pokemonID, versionGroupID); err != nil {
			return fmt.Errorf("failed to sync pokemon %d: %w", pokemonID, err)
		}

		// Create pokedex entry (using array index + 1 as entry number)
		log.Printf("  Inserting pokedex entry for species %d...", pokemonID)
		err = g.pokedexSyncer.InsertPokedexEntry(&external.PokedexEntry{
			PokedexID:   virtualPokedexID,
			SpeciesID:   pokemonID,
			EntryNumber: i + 1, // Sequential numbering
		})
		if err != nil {
			return fmt.Errorf("failed to insert pokedex entry for species %d: %w", pokemonID, err)
		}
	}

	log.Printf("✓ Completed syncing %d Pokemon for special game", len(pokemonIDs))
	return nil
}

// syncPokemonData syncs a single Pokemon including its types, moves, and abilities
// versionGroupID is used to filter which moves to insert (Pokemon learn different moves in different games)
func (g *GameSyncer) syncPokemonData(pokemonID int, versionGroupID int) error {
	// Sync Pokemon
	pokemon, err := g.pokemonSyncer.SyncPokemon(pokemonID)
	if err != nil {
		return fmt.Errorf("failed to sync pokemon: %w", err)
	}

	// Insert types
	log.Printf("    Inserting %d types for pokemon %d (%s)...", len(pokemon.Types), pokemon.ID, pokemon.Name)
	for _, t := range pokemon.Types {
		if err := g.pokemonSyncer.InsertType(&t, pokemon.ID); err != nil {
			return fmt.Errorf("failed to insert type %s for pokemon %d: %w", t.Type.Name, pokemon.ID, err)
		}
	}

	// Sync and insert moves
	log.Printf("    Syncing %d moves for pokemon %d (%s)...", len(pokemon.Moves), pokemon.ID, pokemon.Name)
	movesInserted := 0
	for _, m := range pokemon.Moves {
		moveId, err := utils.ExtractIDFromURL(m.Move.Url)
		if err != nil {
			return fmt.Errorf("failed to extract move ID: %w", err)
		}

		// Sync the move itself (inserts into moves table)
		if err := g.moveSyncer.SyncMove(moveId); err != nil {
			return fmt.Errorf("failed to sync move %d for pokemon %d: %w", moveId, pokemon.ID, err)
		}

		// Insert pokemon_moves for this version group only
		for _, vgDetail := range m.VersionGroupDetails {
			vgID, err := utils.ExtractIDFromURL(vgDetail.VersionGroup.Url)
			if err != nil {
				return fmt.Errorf("failed to extract version group ID: %w", err)
			}

			// Only insert if this move is learned in the current version group
			if vgID == versionGroupID {
				err = g.moveSyncer.repo.InsertPokemonMove(
					pokemon.ID,
					moveId,
					versionGroupID,
					vgDetail.MoveLearnMethod.Name,
					vgDetail.LevelLearnedAt,
				)
				if err != nil {
					return fmt.Errorf("failed to insert pokemon_move: %w", err)
				}
				movesInserted++
			}
		}
	}
	log.Printf("    Inserted %d pokemon_moves for version_group %d", movesInserted, versionGroupID)

	// Insert abilities
	log.Printf("    Inserting %d abilities for pokemon %d (%s)...", len(pokemon.Abilities), pokemon.ID, pokemon.Name)
	for _, a := range pokemon.Abilities {
		if err := g.pokemonSyncer.InsertAbility(&a, pokemon.ID); err != nil {
			return fmt.Errorf("failed to insert ability %s for pokemon %d: %w", a.Ability.Name, pokemon.ID, err)
		}
	}

	return nil
}
