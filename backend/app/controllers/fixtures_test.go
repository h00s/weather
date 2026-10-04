package controllers_test

import (
	"encoding/json"
	"fmt"
	"time"
)

// canned answers each upstream path a service requests. Every task that adds an
// upstream adds its path here; ok is false for a path no service should call.
func canned(path string, now time.Time) (body []byte, ok bool) {
	switch path {
	case "/forecast":
		return forecastJSON(now), true
	}
	return nil, false
}

// forecastJSON is an Open-Meteo forecast for Europe/Zagreb covering the seven
// days from today, so the hourly window always contains now.
func forecastJSON(now time.Time) []byte {
	zagreb, _ := time.LoadLocation("Europe/Zagreb")
	now = now.In(zagreb)
	_, offset := now.Zone()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, zagreb)

	var hourlyTimes []string
	var hourlyTemps []float64
	for i := range 7 * 24 {
		hourlyTimes = append(hourlyTimes, midnight.Add(time.Duration(i)*time.Hour).Format("2006-01-02T15:04"))
		hourlyTemps = append(hourlyTemps, float64(i%24))
	}
	var days, sunrises, sunsets []string
	var maxes []float64
	for i := range 7 {
		day := midnight.AddDate(0, 0, i)
		days = append(days, day.Format(time.DateOnly))
		sunrises = append(sunrises, day.Add(6*time.Hour+46*time.Minute).Format("2006-01-02T15:04"))
		sunsets = append(sunsets, day.Add(18*time.Hour+35*time.Minute).Format("2006-01-02T15:04"))
		maxes = append(maxes, 20+float64(i))
	}

	body, _ := json.Marshal(map[string]any{
		"latitude":              45.59,
		"longitude":             17.23,
		"elevation":             161,
		"utc_offset_seconds":    offset,
		"timezone":              "Europe/Zagreb",
		"timezone_abbreviation": fmt.Sprintf("GMT+%d", offset/3600),
		"current": map[string]any{
			"time": now.Truncate(15 * time.Minute).Format("2006-01-02T15:04"), "interval": 900,
			"temperature_2m": 17.4, "weather_code": 61, "is_day": 1, "wind_direction_10m": 225,
		},
		"hourly": map[string]any{"time": hourlyTimes, "temperature_2m": hourlyTemps},
		"daily":  map[string]any{"time": days, "temperature_2m_max": maxes, "sunrise": sunrises, "sunset": sunsets},
	})
	return body
}
