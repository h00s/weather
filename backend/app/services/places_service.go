package services

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

// placesTSV is written by cmd/genplaces from GeoNames (CC BY 4.0).
//
//go:embed data/hr-places.tsv
var placesTSV string

// PlacesService names locations after Croatian settlements.
type PlacesService struct {
	raptor.Service

	places []models.Place
}

func (s *PlacesService) Setup() error {
	places, err := models.ParsePlaces(strings.NewReader(placesTSV))
	if err != nil {
		return fmt.Errorf("places dataset: %w", err)
	}
	s.places = places
	return nil
}

// Nearest is the settlement at at; ok is false abroad or out at sea.
func (s *PlacesService) Nearest(at models.Coordinates) (models.Place, bool) {
	return models.NearestPlace(s.places, at)
}
