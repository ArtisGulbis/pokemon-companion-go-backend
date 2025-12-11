package db

import (
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	models "github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertVersion(t *testing.T) {
	db := setupTest(t)

	version := &models.Version{
		ID:   1,
		Name: "red",
		Names: []models.Response{
			{
				Name: "yellow",
				Url:  "",
			},
		},
		VersionGroup: models.Response{
			Name: "yellow",
			Url:  "https://pokeapi.co/api/v2/version-group/1/",
		},
	}
	_, err := db.Exec(`INSERT INTO version_groups (id, name, generation_name) VALUES (1, "yellow", "yellow")`)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	repo := NewVersionRepository(db)

	err = repo.InsertVersion(version)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	got, err := repo.GetVersionByID(1)
	require.NoError(t, err)
	require.NotNil(t, got)

	expected := &dto.Version{
		ID:             1,
		Name:           "red",
		DisplayName:    "red",
		VersionGroupID: 1,
	}
	assert.Equal(t, expected, got)
}
