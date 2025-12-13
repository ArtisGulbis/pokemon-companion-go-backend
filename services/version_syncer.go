package services

import (
	"fmt"
	"log"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
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

func (s *VersionSyncer) InsertVersion(v *external.Version) error {
	return s.repo.InsertVersion(v)
}

func (s *VersionSyncer) InsertVersionGroup(vg *external.VersionGroup) error {
	return s.repo.InsertVersionGroup(vg)
}

func (s *VersionSyncer) FetchVersion(id int) (*external.Version, error) {
	return s.client.FetchVersion(id)
}

func (s *VersionSyncer) FetchVersionGroup(id int) (*external.VersionGroup, error) {
	return s.client.FetchVersionGroup(id)
}

func (s *VersionSyncer) SyncVersion(id int) (*external.Version, error) {
	version, err := s.client.FetchVersion(id)
	if err != nil {
		return nil, err
	}
	if err := s.InsertVersion(version); err != nil {
		return nil, err
	}
	return version, nil
}

func (s *VersionSyncer) SyncVersionGroup(id int) (*external.VersionGroup, error) {
	versionGroup, err := s.client.FetchVersionGroup(id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.InsertVersionGroup(versionGroup); err != nil {
		return nil, err
	}

	return versionGroup, nil
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
