package services

import (
	"fmt"
	"log"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/utils"
)

type VersionSyncer struct {
	client      VersionAPIClient
	repo        VersionRepo
	rateLimiter *time.Ticker
}

func NewVersionSyncer(client VersionAPIClient, repo VersionRepo, rateLimiter *time.Ticker) *VersionSyncer {
	return &VersionSyncer{
		client:      client,
		repo:        repo,
		rateLimiter: rateLimiter,
	}
}

func (s *VersionSyncer) SyncAll(limit int) error {
	allVersions, err := s.client.FetchAll(fmt.Sprintf("version?limit=%d", limit))
	if err != nil {
		log.Fatal(err)
	}

	for i, av := range allVersions {
		if i > 0 {
			<-s.rateLimiter.C
		}

		id, err := utils.ExtractIDFromURL(av.Url)
		if err != nil {
			log.Fatal(err)
		}
		version, err := s.client.FetchVersion(id)
		if err != nil {
			log.Fatal(err)
		}
		err = s.repo.InsertVersion(version)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted Version %s (%d/%d)\n", version.Name, i+1, len(allVersions))
	}

	return nil
}
