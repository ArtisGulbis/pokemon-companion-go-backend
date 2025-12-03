package models

type Pokemon struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
}

type Response struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PaginatedResponse struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Response `json:"results"`
}
