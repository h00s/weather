package models

import (
	"time"

	"github.com/h00s/goopenmeteo"
)

// PollenType pairs the name the API returns with its Open-Meteo variable and
// the grains/m³ at which it becomes moderate, high and very high. The bands are
// indicative (common European pollen-count scales) and tuned in this table only.
type PollenType struct {
	Name                     string
	Variable                 goopenmeteo.AirQualityVariable
	Moderate, High, VeryHigh float64
}

var PollenTypes = []PollenType{
	{"alder", goopenmeteo.AlderPollen, 10, 50, 500},
	{"birch", goopenmeteo.BirchPollen, 10, 50, 500},
	{"olive", goopenmeteo.OlivePollen, 20, 100, 500},
	{"grass", goopenmeteo.GrassPollen, 10, 50, 150},
	{"mugwort", goopenmeteo.MugwortPollen, 5, 20, 50},
	{"ragweed", goopenmeteo.RagweedPollen, 5, 20, 50},
}

// AirQualityVariables is what AirQualityService requests, current and hourly.
func AirQualityVariables() goopenmeteo.AirQualityVariables {
	vars := goopenmeteo.AirQualityVariables{goopenmeteo.EuropeanAQI}
	for _, p := range PollenTypes {
		vars = append(vars, p.Variable)
	}
	return vars
}

type AirQualityResponse struct {
	AQI       AQIResponse      `json:"aqi"`
	Pollen    []PollenResponse `json:"pollen"` // every type in PollenTypes, in that order
	FetchedAt time.Time        `json:"fetchedAt"`
}

// AQIResponse is the current European Air Quality Index.
type AQIResponse struct {
	Value float64 `json:"value"`
	Level string  `json:"level"` // good, fair, moderate, poor, very_poor, extremely_poor
}

type PollenResponse struct {
	Type     string  `json:"type"`
	Current  float64 `json:"current"`  // grains/m³; 0 outside the season
	TodayMax float64 `json:"todayMax"` // grains/m³
	Level    string  `json:"level"`    // of TodayMax: none, low, moderate, high, very_high
}

// NewAirQualityResponse maps a one-day air quality forecast.
func NewAirQualityResponse(aq *goopenmeteo.AirQuality, fetchedAt time.Time) AirQualityResponse {
	aqi := aq.Current.Data[goopenmeteo.EuropeanAQI]
	res := AirQualityResponse{
		AQI:       AQIResponse{Value: aqi, Level: AQILevel(aqi)},
		Pollen:    make([]PollenResponse, len(PollenTypes)),
		FetchedAt: fetchedAt,
	}
	for i, p := range PollenTypes {
		todayMax := max(aq.Current.Data[p.Variable], maxOf(aq.Hourly.Data[p.Variable]))
		res.Pollen[i] = PollenResponse{
			Type:     p.Name,
			Current:  aq.Current.Data[p.Variable],
			TodayMax: todayMax,
			Level:    PollenLevel(p.Name, todayMax),
		}
	}
	return res
}

// AQILevel is the European AQI band of value.
func AQILevel(value float64) string {
	switch {
	case value < 20:
		return "good"
	case value < 40:
		return "fair"
	case value < 60:
		return "moderate"
	case value < 80:
		return "poor"
	case value < 100:
		return "very_poor"
	default:
		return "extremely_poor"
	}
}

// PollenLevel is the band of grains/m³ for a pollen type; an unknown type is
// "none".
func PollenLevel(name string, value float64) string {
	for _, p := range PollenTypes {
		if p.Name != name {
			continue
		}
		switch {
		case value >= p.VeryHigh:
			return "very_high"
		case value >= p.High:
			return "high"
		case value >= p.Moderate:
			return "moderate"
		case value >= 1:
			return "low"
		}
	}
	return "none"
}

func maxOf(values []float64) float64 {
	m := 0.0
	for _, v := range values {
		m = max(m, v)
	}
	return m
}
