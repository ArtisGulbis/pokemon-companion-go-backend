package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	SizeCoverSmall     = "cover_small"     // 90×128
	SizeCoverBig       = "cover_big"       // 264×374
	SizeScreenshotMed  = "screenshot_med"  // 569×320
	SizeScreenshotBig  = "screenshot_big"  // 889×500
	SizeScreenshotHuge = "screenshot_huge" // 1280×720
	Size720p           = "720p"            // 1280×720
	Size1080p          = "1080p"           // 1920×1080
)

// Map of Pokemon version names to their IGDB game IDs
var pokemonGameIDs = map[string]int{
	"red":               1561,
	"blue":              1511,
	"yellow":            1512,
	"gold":              1513,
	"silver":            1514,
	"crystal":           1515,
	"ruby":              1516,
	"sapphire":          1517,
	"emerald":           1518,
	"firered":           1519,
	"leafgreen":         1520,
	"diamond":           1521,
	"pearl":             1522,
	"platinum":          1523,
	"heartgold":         1524,
	"soulsilver":        1525,
	"black":             1526,
	"white":             1527,
	"black-2":           1528,
	"white-2":           1529,
	"x":                 7342,
	"y":                 7343,
	"omega-ruby":        8956,
	"alpha-sapphire":    8957,
	"sun":               26757,
	"moon":              26758,
	"ultra-sun":         37489,
	"ultra-moon":        37490,
	"lets-go-pikachu":   105148,
	"lets-go-eevee":     105149,
	"sword":             103577,
	"shield":            103578,
	"brilliant-diamond": 119070,
	"shining-pearl":     119071,
	"legends-arceus":    119388,
	"scarlet":           230614,
	"violet":            230615,
}

type Cover struct {
	ID      int    `json:"id"`
	ImageID string `json:"image_id"`
	URL     string `json:"url"`
}

type Game struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Cover            Cover  `json:"cover"`
	FirstReleaseDate int64  `json:"first_release_date"` // Unix timestamp
}

type IGDBClient struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	TokenExpiry  time.Time
	HTTPClient   *http.Client
}

func NewIGDBClient(clientID, clientSecret string) *IGDBClient {
	return &IGDBClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// Ensure we have a valid access token
func (c *IGDBClient) ensureToken() error {
	if c.AccessToken != "" && time.Now().Before(c.TokenExpiry) {
		return nil
	}

	if c.ClientID == "" || c.ClientSecret == "" {
		return fmt.Errorf("IGDB_CLIENT_ID and IGDB_CLIENT_SECRET must be set")
	}

	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("grant_type", "client_credentials")

	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", data)
	if err != nil {
		return fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OAuth failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return fmt.Errorf("received empty access token")
	}

	c.AccessToken = tokenResp.AccessToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// Get Pokemon game cover by version name using known game IDs
func (c *IGDBClient) GetPokemonGameCover(versionName string) (*Game, error) {
	gameID, exists := pokemonGameIDs[versionName]
	if !exists {
		// Version not in our map, return nil (no cover)
		return nil, nil
	}

	return c.GetGameByID(gameID)
}

// Get a specific game by ID
func (c *IGDBClient) GetGameByID(gameID int) (*Game, error) {
	if err := c.ensureToken(); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`fields name, first_release_date, cover.image_id; where id = %d;`, gameID)

	req, err := http.NewRequest("POST", "https://api.igdb.com/v4/games", bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IGDB API error (status %d): %s", resp.StatusCode, string(body))
	}

	var games []Game
	if err := json.Unmarshal(body, &games); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	if len(games) == 0 {
		return nil, fmt.Errorf("game with ID %d not found", gameID)
	}
	fmt.Println(games[0].FirstReleaseDate)

	return &games[0], nil
}

// Search for a game (generic search - use GetPokemonGameCover for Pokemon games)
func (c *IGDBClient) SearchGame(gameName string) ([]Game, error) {
	if err := c.ensureToken(); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`search "%s"; fields name, cover.image_id; limit 5;`, gameName)

	req, err := http.NewRequest("POST", "https://api.igdb.com/v4/games", bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IGDB API error (status %d): %s", resp.StatusCode, string(body))
	}

	var games []Game
	if err := json.Unmarshal(body, &games); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return games, nil
}

// Get cover URL with specified size
func GetCoverURL(imageID, size string) string {
	if size == "" {
		size = "cover_big"
	}
	return fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_%s/%s.jpg", size, imageID)
}
