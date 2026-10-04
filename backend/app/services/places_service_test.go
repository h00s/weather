package services

import (
	"strings"
	"testing"

	"github.com/h00s/weather/app/models"
)

// The embedded dataset names real points the way a reader would.
func TestEmbeddedPlacesNameRealPoints(t *testing.T) {
	places, err := models.ParsePlaces(strings.NewReader(placesTSV))
	if err != nil {
		t.Fatal(err)
	}
	if len(places) < 10_000 {
		t.Fatalf("%d places, want the whole country", len(places))
	}

	cases := []struct {
		name         string
		at           models.Coordinates
		want, county string
	}{
		{"Zagreb, main square", models.Coordinates{Latitude: 45.813, Longitude: 15.977}, "Zagreb", "21"},
		{"Zagreb, Maksimir", models.Coordinates{Latitude: 45.82, Longitude: 16.02}, "Zagreb", "21"},
		{"Sesvete", models.Coordinates{Latitude: 45.83, Longitude: 16.11}, "Sesvete", "21"},
		{"Velika Gorica", models.Coordinates{Latitude: 45.71, Longitude: 16.07}, "Velika Gorica", "20"},
		{"Daruvar", models.Coordinates{Latitude: 45.59, Longitude: 17.22}, "Daruvar", "01"},
		{"a field near Daruvar", models.Coordinates{Latitude: 45.62, Longitude: 17.30}, "Gornja Vrijeska", "01"},
		{"Split, Riva", models.Coordinates{Latitude: 43.508, Longitude: 16.44}, "Split", "15"},
		{"Hvar", models.Coordinates{Latitude: 43.172, Longitude: 16.442}, "Hvar", "15"},
	}
	for _, c := range cases {
		got, ok := models.NearestPlace(places, c.at)
		if !ok || got.Name != c.want || got.County != c.county {
			t.Errorf("%s: got %q in %s (%v), want %q in %s", c.name, got.Name, got.County, ok, c.want, c.county)
		}
	}

	for name, at := range map[string]models.Coordinates{
		"Sarajevo":          {Latitude: 43.86, Longitude: 18.41},
		"the open Adriatic": {Latitude: 43.3, Longitude: 15.5},
		"Berlin":            {Latitude: 52.52, Longitude: 13.4},
	} {
		if got, ok := models.NearestPlace(places, at); ok {
			t.Errorf("%s was named %q, want no Croatian place", name, got.Name)
		}
	}
}
