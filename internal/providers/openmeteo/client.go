package openmeteo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
	"github.com/SeeknnDestroy/prayertime-cli/internal/providers/httpx"
)

const defaultBaseURL = "https://geocoding-api.open-meteo.com"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *Client) Search(ctx context.Context, query, countryCode string, limit int) ([]app.Location, error) {
	endpoint, err := url.Parse(c.baseURL + "/v1/search")
	if err != nil {
		return nil, app.NewInternalError("failed to build Open-Meteo URL", c.baseURL, "", err)
	}

	params := endpoint.Query()
	params.Set("name", query)
	params.Set("count", fmt.Sprintf("%d", limit))
	params.Set("language", "en")
	params.Set("format", "json")
	if countryCode != "" {
		params.Set("countryCode", countryCode)
	}
	endpoint.RawQuery = params.Encode()

	var payload searchResponse
	if err := httpx.GetJSON(ctx, c.httpClient, endpoint.String(), &payload); err != nil {
		return nil, err
	}

	results := make([]app.Location, 0, len(payload.Results))
	for _, result := range payload.Results {
		results = append(results, app.Location{
			Name:        result.Name,
			Country:     result.Country,
			CountryCode: result.CountryCode,
			Admin1:      result.Admin1,
			Latitude:    result.Latitude,
			Longitude:   result.Longitude,
			Timezone:    result.Timezone,
		})
	}

	return results, nil
}

type searchResponse struct {
	Results []searchResult `json:"results"`
}

type searchResult struct {
	Name        string  `json:"name"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Admin1      string  `json:"admin1"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
}
