package services

type GameSyncer struct {
	versionSyncer *VersionSyncer
	pokedexSyncer *PokedexSyncer
}

func (g *GameSyncer) SyncGame(id int) error {
	// version, err := g.versionSyncer.SyncVersion(id)
	// if err != nil {
	// 	return err
	// }

	// versionGroupId, err := utils.ExtractIDFromURL(version.VersionGroup.Url)
	// if err != nil {
	// 	return err
	// }
	// versionGroup, err := g.versionSyncer.SyncVersionGroup(versionGroupId)
	// if err != nil {
	// 	return err
	// }

	// for _, pdex := range versionGroup.Pokedexes {
	// 	pokedexId, err := utils.ExtractIDFromURL(pdex.Url)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	pokedex, err := g.pokedexSyncer.SyncPokedex(pokedexId)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
