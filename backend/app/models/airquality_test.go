package models

import (
	"testing"
	"time"

	"github.com/h00s/goopenmeteo"
)

func TestAQILevel(t *testing.T) {
	tests := []struct {
		value float64
		want  string
	}{
		{0, "good"}, {19.9, "good"}, {20, "fair"}, {40, "moderate"}, {60, "poor"}, {80, "very_poor"}, {100, "extremely_poor"}, {250, "extremely_poor"},
	}
	for _, tt := range tests {
		if got := AQILevel(tt.value); got != tt.want {
			t.Errorf("AQILevel(%v) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

func TestPollenLevel(t *testing.T) {
	tests := []struct {
		kind  string
		value float64
		want  string
	}{
		{"ragweed", 0, "none"},
		{"ragweed", 0.4, "none"}, // below one grain per m³ is nothing to report
		{"ragweed", 1, "low"},
		{"ragweed", 5, "moderate"},
		{"ragweed", 20, "high"},
		{"ragweed", 50, "very_high"},
		{"birch", 5, "low"},
		{"birch", 10, "moderate"},
		{"birch", 499, "high"},
		{"grass", 150, "very_high"},
	}
	for _, tt := range tests {
		if got := PollenLevel(tt.kind, tt.value); got != tt.want {
			t.Errorf("PollenLevel(%s, %v) = %q, want %q", tt.kind, tt.value, got, tt.want)
		}
	}
}

func TestNewAirQualityResponse(t *testing.T) {
	midnight := time.Date(2026, 9, 29, 0, 0, 0, 0, zagreb)
	ragweed := make([]float64, 24)
	ragweed[15] = 23.5 // the afternoon peak
	aq := &goopenmeteo.AirQuality{
		Current: goopenmeteo.Current{
			Time: midnight.Add(9 * time.Hour),
			Data: map[string]float64{goopenmeteo.EuropeanAQI: 43, goopenmeteo.RagweedPollen: 3.2, goopenmeteo.BirchPollen: 0},
		},
		Hourly: goopenmeteo.TimeseriesData{
			Time: hours(midnight, 24),
			Data: map[string][]float64{
				goopenmeteo.EuropeanAQI:   seq(24, 30),
				goopenmeteo.RagweedPollen: ragweed,
				goopenmeteo.BirchPollen:   make([]float64, 24),
			},
		},
	}
	fetchedAt := midnight.Add(9 * time.Hour)

	res := NewAirQualityResponse(aq, fetchedAt)

	if res.AQI != (AQIResponse{Value: 43, Level: "moderate"}) {
		t.Errorf("aqi = %+v", res.AQI)
	}
	if len(res.Pollen) != len(PollenTypes) {
		t.Fatalf("pollen has %d types, want one per PollenTypes entry (%d)", len(res.Pollen), len(PollenTypes))
	}
	var gotRagweed, gotBirch PollenResponse
	for _, p := range res.Pollen {
		switch p.Type {
		case "ragweed":
			gotRagweed = p
		case "birch":
			gotBirch = p
		}
	}
	if gotRagweed != (PollenResponse{Type: "ragweed", Current: 3.2, TodayMax: 23.5, Level: "high"}) {
		t.Errorf("ragweed = %+v", gotRagweed)
	}
	if gotBirch.Level != "none" {
		t.Errorf("birch = %+v", gotBirch)
	}
	if !res.FetchedAt.Equal(fetchedAt) {
		t.Errorf("fetchedAt = %v", res.FetchedAt)
	}
}
