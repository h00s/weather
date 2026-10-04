package models

import (
	"fmt"
	"math"
	"strconv"
)

// Coordinates is a point on Earth in decimal degrees.
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

// ParseCoordinates reads lat and lon as sent in a query string. Both must be
// finite numbers in range; the error names the one that isn't.
func ParseCoordinates(lat, lon string) (Coordinates, error) {
	latitude, err := parseDegrees(lat, 90)
	if err != nil {
		return Coordinates{}, fmt.Errorf("lat %q: %w", lat, err)
	}
	longitude, err := parseDegrees(lon, 180)
	if err != nil {
		return Coordinates{}, fmt.Errorf("lon %q: %w", lon, err)
	}
	return Coordinates{latitude, longitude}, nil
}

func parseDegrees(s string, limit float64) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(v) || math.Abs(v) > limit {
		return 0, fmt.Errorf("want a number from -%g to %g", limit, limit)
	}
	return v, nil
}

// Rounded snaps c to 2 decimals, a cell of about 1 km. Every upstream request
// and cache key uses it, so neighbours share one cached answer and no precise
// position is passed on.
func (c Coordinates) Rounded() Coordinates {
	return Coordinates{round2(c.Latitude), round2(c.Longitude)}
}

// Key names c's cell, such as "45.59,17.23".
func (c Coordinates) Key() string {
	r := c.Rounded()
	return strconv.FormatFloat(r.Latitude, 'f', 2, 64) + "," + strconv.FormatFloat(r.Longitude, 'f', 2, 64)
}

func round2(v float64) float64 {
	r := math.Round(v*100) / 100
	if r == 0 {
		return 0 // never "-0.00"
	}
	return r
}
