package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadImage downloads an image from a URL and saves it to the specified path
func DownloadImage(url, savePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(savePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file already exists
	if _, err := os.Stat(savePath); err == nil {
		// File exists, skip download
		return nil
	}

	// Download the image
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
}

// GetPokemonSpritePath returns the local path for a Pokemon sprite
func GetPokemonSpritePath(pokemonID int, spriteType string) string {
	return filepath.Join("images", "pokemon", fmt.Sprintf("%d_%s.png", pokemonID, spriteType))
}

// GetGameCoverPath returns the local path for a game cover
func GetGameCoverPath(versionName string) string {
	return filepath.Join("images", "covers", fmt.Sprintf("%s.jpg", versionName))
}
