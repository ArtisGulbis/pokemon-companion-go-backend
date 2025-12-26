package db

import (
	"testing"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/dto"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertMove(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewMoveRepository(db)

	move := &external.Move{
		ID:   1,
		Name: "tackle",
		Type: external.Response{
			Name: "normal",
			Url:  "",
		},
		Power:    100,
		Accuracy: 100,
		PP:       10,
		DamageClass: external.Response{
			Name: "physical",
			Url:  "",
		},
		EffectEntries: []external.EffectEntry{
			{
				Effect: "",
				Language: external.Response{
					Name: "en",
					Url:  "",
				},
				ShortEffect: "hit the enemy",
			},
		},
		Priority: 0,
	}

	err = repo.InsertMove(move)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := repo.GetMoveByID(1)
	require.NoError(t, err)
	require.NotNil(t, actual)

	expected := &dto.Move{
		ID:          1,
		Name:        "tackle",
		Type:        "normal",
		Power:       100,
		Accuracy:    100,
		PP:          10,
		DamageClass: "physical",
		EffectShort: "hit the enemy",
		Priority:    0,
	}

	assert.Equal(t, expected, actual)
}
