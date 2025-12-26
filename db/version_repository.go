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
	versionGroupId, err := utils.ExtractIDFromURL(v.VersionGroup.Url)
	if err != nil {
		return err
	}

	fmt.Printf("DEBUG Version Insert - ID: %d, Name: %s, Cover: %s, ReleaseDate: %d, VersionGroupID: %d\n",
		v.ID, v.Name, v.Cover, v.ReleaseDate, versionGroupId)

	// Verify version_group exists before inserting
	var vgExists int
	checkErr := r.db.QueryRow("SELECT COUNT(*) FROM version_groups WHERE id = ?", versionGroupId).Scan(&vgExists)
	if checkErr != nil {
		fmt.Printf("DEBUG: Error checking version_group: %v\n", checkErr)
	} else {
		fmt.Printf("DEBUG: Version_group %d exists before version insert: %v\n", versionGroupId, vgExists > 0)
	}

	// Use direct SQL instead of prepared statement
	directSQL := "INSERT OR IGNORE INTO versions (id, name, cover, release_date, display_name, version_group_id) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := r.db.Exec(directSQL, v.ID, v.Name, v.Cover, v.ReleaseDate, v.Name, versionGroupId)
	if err != nil {
		return fmt.Errorf("version insert failed: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("DEBUG: Version insert - RowsAffected: %d\n", rowsAffected)

	return nil
}

func (r *VersionRepository) InsertVersionGroup(v *external.VersionGroup) error {
	stmt, err := r.db.Prepare(queries.InsertVersionGroup)
	if err != nil {
		return err
	}
	defer stmt.Close()

	fmt.Printf("DEBUG: Inserting version_group - ID: %d, Name: %s, Generation: %s\n",
		v.ID, v.Name, v.Generation.Name)

	result, err := stmt.Exec(
		v.ID,
		v.Name,
		v.Generation.Name,
	)
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("DEBUG: RowsAffected: %d\n", rowsAffected)

	// Check if it was actually inserted
	var count int
	checkErr := r.db.QueryRow("SELECT COUNT(*) FROM version_groups WHERE id = ?", v.ID).Scan(&count)
	if checkErr != nil {
		fmt.Printf("DEBUG: Error checking: %v\n", checkErr)
	} else {
		fmt.Printf("DEBUG: Version group %d exists in DB: %v\n", v.ID, count > 0)
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
