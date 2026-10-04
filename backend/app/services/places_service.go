package services

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

// placesTSV is written by cmd/genplaces from GeoNames (CC BY 4.0).
//
//go:embed data/hr-places.tsv
var placesTSV string

const (
	geocodingURL      = "https://geocoding-api.open-meteo.com/v1"
	searchTTL         = 24 * time.Hour // place names don't move
	maxCachedSearches = 2000
	searchLimit       = 10
)

// PlacesService names locations after Croatian settlements and finds them by
// name through Open-Meteo's geocoding API. APP_GEOCODING_URL overrides its root
// (tests).
type PlacesService struct {
	raptor.Service

	Upstream *UpstreamService

	places   []models.Place
	client   *http.Client
	baseURL  string
	searches *keyedCache[string, []models.PlaceResponse]
}

func (s *PlacesService) Setup() error {
	places, err := models.ParsePlaces(strings.NewReader(placesTSV))
	if err != nil {
		return fmt.Errorf("places dataset: %w", err)
	}
	s.places = places
	s.client = &http.Client{Timeout: 10 * time.Second}
	s.baseURL = strings.TrimSuffix(s.Config.AppString("geocoding_url", geocodingURL), "/")
	s.searches = newKeyedCache[string, []models.PlaceResponse](searchTTL, maxCachedSearches)
	return nil
}

// Nearest is the settlement at at; ok is false abroad or out at sea.
func (s *PlacesService) Nearest(at models.Coordinates) (models.Place, bool) {
	return models.NearestPlace(s.places, at)
}

// Search finds Croatian settlements by name, the best match first (Open-Meteo
// matches two letters exactly and three or more fuzzily). When the geocoding
// API fails it returns the last answer to the same query, or the error when
// there is none.
func (s *PlacesService) Search(ctx context.Context, query string) ([]models.PlaceResponse, error) {
	places, _, err := s.searches.Get(ctx, strings.ToLower(query), func(ctx context.Context) ([]models.PlaceResponse, error) {
		places, err := s.geocode(ctx, query)
		if err != nil {
			s.Log.Warn("Open-Meteo geocoding request failed", "error", err)
		}
		return places, err
	})
	return places, err
}

func (s *PlacesService) geocode(ctx context.Context, query string) ([]models.PlaceResponse, error) {
	if err := s.Upstream.Take(); err != nil {
		return nil, err
	}
	params := url.Values{
		"name":        {query},
		"count":       {"20"}, // the filter drops airfields, stations and other features
		"language":    {"hr"},
		"countryCode": {"HR"},
		"format":      {"json"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/search?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding: unexpected status %d", resp.StatusCode)
	}
	var g models.GeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, fmt.Errorf("geocoding: %w", err)
	}
	return models.SearchResults(g, searchLimit), nil
}
