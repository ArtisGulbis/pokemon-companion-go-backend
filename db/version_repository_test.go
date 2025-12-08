package db

import (
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	models "github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertVersion(t *testing.T) {
	t.Skip("Version repository not implemented yet")
	db := setupTest(t)

	version := &models.Version{
		ID:   1,
		Name: "red",
		Names: []models.Response{
			{
				Name: "",
				Url:  "",
			},
		},
		VersionGroup: models.Response{
			Name: "",
			Url:  "",
		},
	}

	repo := NewVersionRepository(db)

	err := repo.InsertVersion(version)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	got, err := repo.GetVersionByID(1)
	require.NoError(t, err)
	require.NotNil(t, got)

	expected := &dto.Version{}
	assert.Equal(t, expected, got)
}
