package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/models"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) FetchAllPokemon() ([]models.Response, error) {
	url := fmt.Sprintf("%s/api/v2/pokemon?limit=1500", c.BaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch all pokemon: %s", resp.Status)
	}
	var allPokemon models.PaginatedResponse

	if err := json.NewDecoder(resp.Body).Decode(&allPokemon); err != nil {
		return nil, fmt.Errorf("failed to decode all pokemon :%w", err)
	}

	return allPokemon.Results, nil
}

func (c *Client) FetchPokemon(id int) (*models.Pokemon, error) {
	url := fmt.Sprintf("%s/api/v2/pokemon/%d", c.BaseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch pokemon: %s", resp.Status)
	}

	var pokemon models.Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, fmt.Errorf("failed to decode pokemon: %w", err)
	}

	return &pokemon, nil
}
