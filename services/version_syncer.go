package services

import (
	"fmt"
	"log"
	"time"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/igdb"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
	"github.com/ArtisGulbis/pokemon-companion-go-backend/utils"
)

type VersionSyncer struct {
	client      VersionAPIClient
	igdbClient  *igdb.IGDBClient
	repo        VersionRepo
	rateLimiter *time.Ticker
}

func NewVersionSyncer(client VersionAPIClient, igdbClient *igdb.IGDBClient, repo VersionRepo, rateLimiter *time.Ticker) *VersionSyncer {
	return &VersionSyncer{
		client:      client,
		igdbClient:  igdbClient,
		repo:        repo,
		rateLimiter: rateLimiter,
	}
}

func (s *VersionSyncer) InsertVersion(v *external.Version) error {
	// Try to get cover using the known Pokemon game mapping
	game, err := s.igdbClient.GetPokemonGameCover(v.Name)
	if err != nil {
		return err
	}
	if game != nil && game.Cover.ImageID != "" {
		// Get the cover URL from IGDB
		coverURL := igdb.GetCoverURL(game.Cover.ImageID, "cover_big")

		// Download and save the image locally
		localPath := utils.GetGameCoverPath(v.Name)
		if err := utils.DownloadImage(coverURL, localPath); err != nil {
			log.Printf("Warning: failed to download cover for %s: %v", v.Name, err)
		} else {
			// Store the local path instead of the URL
			v.Cover = localPath
		}
	}
	v.ReleaseDate = int(game.FirstReleaseDate)
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
