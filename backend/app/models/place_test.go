package models

import (
	"strings"
	"testing"
)

func TestParsePlaces(t *testing.T) {
	places, err := ParsePlaces(strings.NewReader("# comment\nDaruvar\t45.59056\t17.225\t01\t7440\n\nZagreb\t45.81444\t15.97798\t21\t663592\n"))
	if err != nil {
		t.Fatal(err)
	}
	want := Place{Name: "Daruvar", Coordinates: Coordinates{45.59056, 17.225}, County: "01", Population: 7440}
	if len(places) != 2 || places[0] != want || places[1].Name != "Zagreb" {
		t.Errorf("places = %+v", places)
	}

	for _, bad := range []string{"Daruvar\t45.59\n", "Daruvar\tx\t17.2\t01\t1\n", "Daruvar\t45.5\t17.2\t01\tmany\n"} {
		if _, err := ParsePlaces(strings.NewReader(bad)); err == nil {
			t.Errorf("ParsePlaces(%q) accepted a bad line", bad)
		}
	}
}

func TestDistanceKm(t *testing.T) {
	zagreb, split := Coordinates{45.81444, 15.97798}, Coordinates{43.50891, 16.43915}
	if d := DistanceKm(zagreb, split); d < 255 || d > 262 {
		t.Errorf("Zagreb–Split = %.1f km, want about 259", d)
	}
	if d := DistanceKm(zagreb, zagreb); d != 0 {
		t.Errorf("distance to itself = %v", d)
	}
}

func TestNearestPlaceLetsLargerPlacesReachFurther(t *testing.T) {
	// One degree of longitude at 45°N is 78.7 km: the village sits 3.9 km east of the city.
	places := []Place{
		{Name: "City", Coordinates: Coordinates{45.0, 16.0}, County: "21", Population: 600_000}, // radius ≈ 14.7 km
		{Name: "Village", Coordinates: Coordinates{45.0, 16.05}, County: "20"},                  // radius 1 km
	}
	cases := []struct {
		at   Coordinates
		want string // "" for no place
	}{
		{Coordinates{45.0, 16.04}, "City"},     // 0.8 km from the village, but well inside the city's reach
		{Coordinates{45.0, 16.049}, "Village"}, // right at the village
		{Coordinates{45.0, 16.2}, "Village"},   // the city is 15.7 km away, beyond reach; the village 11.8 km
		{Coordinates{45.0, 16.3}, ""},          // nothing within 15 km
		{Coordinates{46.5, 16.0}, ""},
	}
	for _, c := range cases {
		got, ok := NearestPlace(places, c.at)
		if c.want == "" && ok {
			t.Errorf("NearestPlace(%v) = %s, want none", c.at, got.Name)
		}
		if c.want != "" && (!ok || got.Name != c.want) {
			t.Errorf("NearestPlace(%v) = %s (%v), want %s", c.at, got.Name, ok, c.want)
		}
	}
}

func TestNewPlaceResponseNamesTheCounty(t *testing.T) {
	got := NewPlaceResponse(Place{Name: "Daruvar", Coordinates: Coordinates{45.59056, 17.225}, County: "01"})
	want := PlaceResponse{Name: "Daruvar", County: "Bjelovarsko-bilogorska županija", Latitude: 45.59056, Longitude: 17.225}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
