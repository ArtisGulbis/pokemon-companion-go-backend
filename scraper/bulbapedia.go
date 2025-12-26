package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
}

func NewScraper() *Scraper {
	return &Scraper{}
}

type GymLeaders struct {
	name string
	url  string
}

func (s *Scraper) ScrapeGamePage(url string) error {
	gymLeaders := findGymLeaders(url)
	fmt.Println(gymLeaders)

	return nil
}

func findGymLeaders(url string) []GymLeaders {
	res, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil
	}

	// Find the "Gyms" heading
	gymsHeading := doc.Find("span#Gyms")
	if gymsHeading.Length() == 0 {
		fmt.Println("Could not find Gyms section")
		return nil
	}

	// Get the parent heading, then find the next <p> tag
	parent := gymsHeading.Parent()
	nextParagraph := parent.NextFilteredUntil("p", "h2, h3, h4")

	if nextParagraph.Length() == 0 {
		fmt.Println("Could not find paragraph with gym leaders")
		return nil
	}

	// Extract gym leader names
	var gymLeaders []GymLeaders
	nextParagraph.Find("a").Each(func(i int, link *goquery.Selection) {
		href, exists := link.Attr("href")
		if !exists {
			return
		}

		// Skip type links (they contain "(type)" in href)
		if len(href) > 0 && !strings.Contains(href, "(type)") {
			name := link.Text()
			if len(name) > 0 {
				gymLeaders = append(gymLeaders, GymLeaders{
					name: name,
					url:  href,
				})
			}
		}
	})

	return gymLeaders
}
