package main

import (
	"log"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/scraper"
)

func main() {
	s := scraper.NewScraper()

	// Test with Pokemon Red/Blue page
	url := "https://bulbapedia.bulbagarden.net/wiki/Pokémon_Red_and_Blue_Versions"

	if err := s.ScrapeGamePage(url); err != nil {
		log.Fatal(err)
	}
}
