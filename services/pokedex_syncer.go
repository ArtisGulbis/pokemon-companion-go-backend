package services

import (
	"fmt"
	"log"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/utils"
)

type PokedexSyncer struct {
	client      PokedexAPIClient
	repo        PokedexRepo
	rateLimiter *time.Ticker
}

func NewPokedexSyncer(client PokedexAPIClient, repo PokedexRepo, rateLimiter *time.Ticker) *PokedexSyncer {
	return &PokedexSyncer{
		client:      client,
		repo:        repo,
		rateLimiter: rateLimiter,
	}
}

func (s *PokedexSyncer) InsertPokedex(pd *external.Pokedex) error {
	return s.repo.InsertPokedex(pd)
}

func (s *PokedexSyncer) InsertVersionGroupPokedex(vg *external.VersionGroup) error {
	return s.repo.InsertVersionGroupPokedex(vg)
}

func (s *PokedexSyncer) FetchPokedex(id int) (*external.Pokedex, error) {
	return s.client.FetchPokedex(id)
}

func (s *PokedexSyncer) SyncPokedex(id int) (*external.Pokedex, error) {
	pokedex, err := s.client.FetchPokedex(id)
	if err != nil {
		return nil, err
	}

	err = s.repo.InsertPokedex(pokedex)
	if err != nil {
		return nil, err
	}

	return pokedex, nil
}

func (s *PokedexSyncer) SyncAll(limit int) error {
	allPokedexes, err := s.client.FetchAll(fmt.Sprintf("pokedex?limit=%d", limit))
	if err != nil {
		log.Fatal(err)
	}

	for i, pkdx := range allPokedexes {
		if i > 0 {
			<-s.rateLimiter.C
		}

		id, err := utils.ExtractIDFromURL(pkdx.Url)
		if err != nil {
			log.Fatal(err)
		}

		pokedex, err := s.client.FetchPokedex(id)
		if err != nil {
			log.Fatal(err)
		}

		err = s.repo.InsertPokedex(pokedex)
		if err != nil {
			log.Fatalf("Failed to insert pokedex ID %d: %v", pokedex.ID, err)
		}

		fmt.Printf("Inserted Pokedex %s (%d/%d)\n", pokedex.Name, i+1, len(allPokedexes))
	}
	return nil
}
