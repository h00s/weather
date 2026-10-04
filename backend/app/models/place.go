package models

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// Place is a Croatian settlement from the GeoNames dataset.
type Place struct {
	Name string
	Coordinates
	County     string // GeoNames admin1 code, a key of Counties
	Population int
}

// PlaceResponse is a named place to show weather for.
type PlaceResponse struct {
	Name      string  `json:"name"`
	County    string  `json:"county"` // e.g. "Bjelovarsko-bilogorska županija"; "" when unknown
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewPlaceResponse(p Place) PlaceResponse {
	return PlaceResponse{Name: p.Name, County: Counties[p.County].Name, Latitude: p.Latitude, Longitude: p.Longitude}
}

// ParsePlaces reads the tab-separated dataset cmd/genplaces writes: name,
// latitude, longitude, county code and population per line; # starts a comment.
func ParsePlaces(r io.Reader) ([]Place, error) {
	var places []Place
	scanner := bufio.NewScanner(r)
	for n := 1; scanner.Scan(); n++ {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) != 5 {
			return nil, fmt.Errorf("line %d: %d fields, want 5", n, len(f))
		}
		lat, errLat := strconv.ParseFloat(f[1], 64)
		lon, errLon := strconv.ParseFloat(f[2], 64)
		population, errPop := strconv.Atoi(f[4])
		if err := cmp.Or(errLat, errLon, errPop); err != nil {
			return nil, fmt.Errorf("line %d: %w", n, err)
		}
		places = append(places, Place{Name: f[0], Coordinates: Coordinates{lat, lon}, County: f[3], Population: population})
	}
	return places, scanner.Err()
}

// placeReachKm is how far a point may be from every settlement and still be
// named: farther, it is abroad or out at sea.
const placeReachKm = 15

// NearestPlace names the settlement a point is in. Larger places reach further:
// distance counts in units of a place's radius, max(1 km, 0.6 km·√(population/1000)),
// so a point in Maksimir is Zagreb (3 km from a city of 660,000) and not the
// village 1.5 km away, while a field near Daruvar is the hamlet beside it.
// ok is false when no settlement lies within placeReachKm.
func NearestPlace(places []Place, at Coordinates) (place Place, ok bool) {
	best := math.Inf(1)
	for _, p := range places {
		// A cheap box first: 0.2° of latitude and 0.3° of longitude both exceed
		// 15 km across Croatia, so this only skips places out of reach.
		if math.Abs(p.Latitude-at.Latitude) > 0.2 || math.Abs(p.Longitude-at.Longitude) > 0.3 {
			continue
		}
		d := DistanceKm(at, p.Coordinates)
		if d > placeReachKm {
			continue
		}
		if score := d / p.radiusKm(); score < best {
			place, best, ok = p, score, true
		}
	}
	return place, ok
}

func (p Place) radiusKm() float64 {
	return max(1, 0.6*math.Sqrt(float64(p.Population)/1000))
}

// DistanceKm is the great-circle distance between a and b.
func DistanceKm(a, b Coordinates) float64 {
	const earthRadiusKm = 6371
	rad := math.Pi / 180
	dLat := (b.Latitude - a.Latitude) * rad
	dLon := (b.Longitude - a.Longitude) * rad
	h := math.Pow(math.Sin(dLat/2), 2) + math.Cos(a.Latitude*rad)*math.Cos(b.Latitude*rad)*math.Pow(math.Sin(dLon/2), 2)
	return 2 * earthRadiusKm * math.Asin(math.Sqrt(h))
}
