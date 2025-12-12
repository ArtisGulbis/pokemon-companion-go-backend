package db

import (
	"fmt"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/queries"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/utils"
)

type VersionRepository struct {
	db *Database
}

func NewVersionRepository(db *Database) *VersionRepository {
	return &VersionRepository{db: db}
}

func (r *VersionRepository) InsertVersion(v *external.Version) error {
	stmt, err := r.db.Prepare(queries.InsertVersion)
	if err != nil {
		return err
	}
	defer stmt.Close()

	versionGroupId, err := utils.ExtractIDFromURL(v.VersionGroup.Url)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		v.ID,
		v.Name,
		v.Name,
		versionGroupId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *VersionRepository) InsertVersionGroup(v *external.VersionGroup) error {
	stmt, err := r.db.Prepare(queries.InsertVersionGroup)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		v.ID,
		v.Name,
		v.Generation.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *VersionRepository) GetVersionByID(id int) (*dto.Version, error) {
	rows, err := r.db.Query(queries.GetVersionByID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var version *dto.Version

	if rows.Next() {
		version = &dto.Version{}
		err = rows.Scan(
			&version.ID,
			&version.Name,
			&version.DisplayName,
			&version.VersionGroupID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if version == nil {
		return nil, fmt.Errorf("version %d not found", id)
	}

	return version, nil
}
