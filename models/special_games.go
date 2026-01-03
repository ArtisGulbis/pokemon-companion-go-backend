package models

// Special game versions that don't have traditional Pokedexes in PokeAPI
// These require hardcoded Pokemon lists

// ColosseumPokemonIDs contains all Pokemon available in Pokemon Colosseum
var ColosseumPokemonIDs = []int{
	153, 156, 159, 162, 164, 166, 168, 176, 180, 185,
	188, 190, 192, 193, 195, 196, 197, 198, 200, 205,
	206, 207, 210, 211, 213, 214, 215, 217, 218, 221,
	223, 225, 226, 227, 229, 234, 235, 237, 241, 243,
	244, 245, 248, 296, 307, 329, 333, 357, 359, 376,
}

// XDPokemonIDs contains all Pokemon available in Pokemon XD: Gale of Darkness
var XDPokemonIDs = []int{
	133, 216, 165, 261, 228, 343, 363, 179, 316, 273,
	167, 322, 285, 318, 301, 100, 296, 37, 355, 280,
	303, 361, 204, 177, 315, 52, 220, 21, 88, 337,
	299, 335, 46, 58, 90, 15, 17, 114, 12, 82, 49,
	70, 24, 57, 97, 55, 302, 85, 20, 83, 334, 115,
	354, 126, 127, 78, 219, 107, 106, 108, 123,
	113, 338, 121, 277, 125, 143, 62, 122, 51,
	310, 373, 105, 131, 112, 103, 128, 149,
}

// GetSpecialGamePokemon returns the Pokemon list for special game versions
// Returns nil if the version doesn't have a special list
func GetSpecialGamePokemon(versionName string) []int {
	switch versionName {
	case "colosseum":
		return ColosseumPokemonIDs
	case "xd":
		return XDPokemonIDs
	default:
		return nil
	}
}
