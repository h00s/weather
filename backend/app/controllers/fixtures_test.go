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
	case "/search":
		return []byte(searchJSON), true
	case "/api/v1/warnings/feeds-croatia":
		return meteoalarmJSON(now), true
	case "/air-quality":
		return airQualityJSON(now), true
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

// searchJSON is Open-Meteo's geocoding answer for "Daru": two settlements and an airfield.
const searchJSON = `{"results":[
 {"id":3202184,"name":"Daruvar","latitude":45.59056,"longitude":17.225,"feature_code":"PPLA2","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija"},
 {"id":11500092,"name":"Daruvar","latitude":45.58507,"longitude":17.2114,"feature_code":"AIRF","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija"},
 {"id":12509853,"name":"Daruvarski Vinogradi","latitude":45.60251,"longitude":17.25084,"feature_code":"PPL","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija"}
],"generationtime_ms":0.3}`

// meteoalarmJSON is DHMZ's Meteoalarm feed with a yellow wind warning in force
// for Bjelovarsko-bilogorska and an orange one for Zadarska.
func meteoalarmJSON(now time.Time) []byte {
	warning := func(id, area, level string) map[string]any {
		return map[string]any{"alert": map[string]any{
			"identifier": id, "msgType": "Alert",
			"info": []map[string]any{{
				"language": "hr-HR", "event": "Upozorenje", "description": "Opis",
				"onset": now.Add(-time.Hour).Format(time.RFC3339), "expires": now.Add(5 * time.Hour).Format(time.RFC3339),
				"responseType": []string{"Monitor"},
				"parameter": []map[string]string{
					{"valueName": "awareness_level", "value": level},
					{"valueName": "awareness_type", "value": "1; Wind"},
				},
				"area": []map[string]any{{"areaDesc": area, "geocode": []map[string]string{{"valueName": "EMMA_ID", "value": "X"}}}},
			}},
		}}
	}
	body, _ := json.Marshal(map[string]any{"warnings": []any{
		warning("here", "Bjelovarsko-bilogorska", "2; yellow; Moderate"),
		warning("coast", "Zadarska", "3; orange; Severe"),
	}})
	return body
}

// airQualityJSON is an Open-Meteo air quality answer for today: moderate AQI and
// a ragweed peak of 23.5 grains/m³ in the afternoon.
func airQualityJSON(now time.Time) []byte {
	zagreb, _ := time.LoadLocation("Europe/Zagreb")
	now = now.In(zagreb)
	_, offset := now.Zone()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, zagreb)
	var times []string
	ragweed := make([]float64, 24)
	for i := range 24 {
		times = append(times, midnight.Add(time.Duration(i)*time.Hour).Format("2006-01-02T15:04"))
	}
	ragweed[15] = 23.5

	body, _ := json.Marshal(map[string]any{
		"utc_offset_seconds":    offset,
		"timezone":              "Europe/Zagreb",
		"timezone_abbreviation": fmt.Sprintf("GMT+%d", offset/3600),
		"current": map[string]any{
			"time": now.Truncate(time.Hour).Format("2006-01-02T15:04"), "interval": 3600,
			"european_aqi": 43, "ragweed_pollen": 3.2, "birch_pollen": nil,
		},
		"hourly": map[string]any{"time": times, "ragweed_pollen": ragweed},
	})
	return body
}
