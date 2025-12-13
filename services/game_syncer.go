package services

import "github.com/ArtisGulbis/pokemon-companion-go-backend/utils"

type GameSyncer struct {
	versionSyncer *VersionSyncer
	pokedexSyncer *PokedexSyncer
	pokemonSyncer *PokemonSyncer
}

func NewGameSyncer(
	versionSyncer *VersionSyncer,
	pokedexSyncer *PokedexSyncer,
	pokemonSyncer *PokemonSyncer,
) *GameSyncer {
	return &GameSyncer{
		versionSyncer: versionSyncer,
		pokedexSyncer: pokedexSyncer,
		pokemonSyncer: pokemonSyncer,
	}
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
	}

	return nil
}
