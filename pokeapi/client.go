package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/ArtisGulbis/pokemon-companion-go-backend/models/external"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) FetchPokemon(id int) (*models.Pokemon, error) {
	return fetchByID[models.Pokemon](c, "pokemon", id)
}

func (c *Client) FetchVersion(id int) (*models.Version, error) {
	return fetchByID[models.Version](c, "version", id)
}

func (c *Client) FetchVersionGroup(id int) (*models.VersionGroup, error) {
	return fetchByID[models.VersionGroup](c, "version-group", id)
}

func (c *Client) FetchPokedex(id int) (*models.Pokedex, error) {
	return fetchByID[models.Pokedex](c, "pokedex", id)
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

func fetchByID[T any](c *Client, resource string, id int) (*T, error) {
	url := fmt.Sprintf("%s/api/v2/%s/%d", c.BaseURL, resource, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s: %s", resource, resp.Status)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", resource, err)
	}

	return &result, nil
}
