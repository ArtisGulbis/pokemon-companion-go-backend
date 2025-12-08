package db

import (
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
)

type VersionRepository struct {
	db *Database
}

func NewVersionRepository(db *Database) *VersionRepository {
	return &VersionRepository{db: db}
}

func (r *VersionRepository) InsertVersion(v *external.Version) error {
	//stmt, err := r.db.Prepare(queries.InsertVersion)
	return nil
}

func (r *VersionRepository) GetVersionByID(id int) (*dto.Version, error) {
	return nil, nil
}
