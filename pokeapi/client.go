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

func (c *Client) FetchPokedex(id int) (*models.Pokedex, error) {
	url := fmt.Sprintf("%s/api/v2/pokedex/%d", c.BaseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch pokedex: %s", resp.Status)
	}

	var pokedex models.Pokedex
	if err := json.NewDecoder(resp.Body).Decode(&pokedex); err != nil {
		return nil, fmt.Errorf("failed to decode pokedex: %w", err)
	}

	return &pokedex, nil
}

func (c *Client) FetchAll(path string) ([]models.Response, error) {
	url := fmt.Sprintf("%s/api/v2/%s", c.BaseURL, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch all pokemon: %s", resp.Status)
	}
	var paginatedResponse models.PaginatedResponse

	if err := json.NewDecoder(resp.Body).Decode(&paginatedResponse); err != nil {
		return nil, fmt.Errorf("failed to decode all pokemon :%w", err)
	}

	return paginatedResponse.Results, nil
}
