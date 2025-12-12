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

	// versionGroup, err := g.versionSyncer.SyncVersionGroup(version.VersionGroup.Name)
	// if err != nil {
	// 	return err
	// }

	// for _, pdex := range versionGroup.Pokedexes {
	// 	pokedex := g.pokedexSyncer.SyncPokedex(pdex.Name)
	// }

	return nil
}
