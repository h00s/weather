package models

import "strings"

// GeocodingResponse is Open-Meteo's geocoding answer
// (https://open-meteo.com/en/docs/geocoding-api), reduced to what search uses.
// results is absent when nothing matches.
type GeocodingResponse struct {
	Results []GeocodingResult `json:"results"`
}

type GeocodingResult struct {
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	FeatureCode string  `json:"feature_code"` // GeoNames feature code: PPL… for settlements
	CountryCode string  `json:"country_code"`
	Admin1      string  `json:"admin1"` // the county, in Croatian with language=hr
}

// gone are settlements that no longer exist: historical, abandoned, destroyed.
var gone = map[string]bool{"PPLH": true, "PPLQ": true, "PPLW": true, "PPLCH": true}

// SearchResults keeps the Croatian settlements, at most limit, in Open-Meteo's
// order: the best match first, larger places before smaller ones of one name.
func SearchResults(g GeocodingResponse, limit int) []PlaceResponse {
	places := make([]PlaceResponse, 0, min(limit, len(g.Results)))
	for _, r := range g.Results {
		if len(places) == limit {
			break
		}
		if r.CountryCode != "HR" || !strings.HasPrefix(r.FeatureCode, "PPL") || gone[r.FeatureCode] {
			continue
		}
		places = append(places, PlaceResponse{Name: r.Name, County: r.Admin1, Latitude: r.Latitude, Longitude: r.Longitude})
	}
	return places
}
