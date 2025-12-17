package dto

type Pokemon struct {
	ID                 int
	SpeciesID          int
	Name               string
	IsDefault          bool
	Height             int
	Weight             int
	BaseExperience     int
	HP                 int
	Attack             int
	Defense            int
	SpecialAttack      int
	SpecialDefense     int
	Speed              int
	SpriteFrontDefault string
	SpriteFrontShiny   string
	SpriteArtwork      string
}
