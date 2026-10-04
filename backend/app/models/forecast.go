package models

import (
	"time"

	"github.com/h00s/goopenmeteo"
)

// The Open-Meteo variables NewForecastResponse reads; ForecastService requests
// exactly these.
var (
	ForecastCurrentVariables = goopenmeteo.WeatherVariables{
		goopenmeteo.Temperature2M,
		goopenmeteo.ApparentTemperature,
		goopenmeteo.RelativeHumidity2M,
		goopenmeteo.DewPoint2M,
		goopenmeteo.Precipitation,
		goopenmeteo.WeatherCode,
		goopenmeteo.CloudCover,
		goopenmeteo.PressureMsl,
		goopenmeteo.WindSpeed10M,
		goopenmeteo.WindDirection10M,
		goopenmeteo.WindGusts10M,
		goopenmeteo.Visibility,
		goopenmeteo.UVIndex, // declared with the air-quality variables; the forecast API has it too
		goopenmeteo.IsDay,
	}
	ForecastHourlyVariables = goopenmeteo.WeatherVariables{
		goopenmeteo.Temperature2M,
		goopenmeteo.ApparentTemperature,
		goopenmeteo.PrecipitationProbability,
		goopenmeteo.Precipitation,
		goopenmeteo.WeatherCode,
		goopenmeteo.WindSpeed10M,
		goopenmeteo.WindDirection10M,
		goopenmeteo.IsDay,
	}
	ForecastDailyVariables = goopenmeteo.WeatherVariables{
		goopenmeteo.Temperature2MMin,
		goopenmeteo.Temperature2MMax,
		goopenmeteo.WeatherCode,
		goopenmeteo.PrecipitationSum,
		goopenmeteo.PrecipitationProbabilityMax,
		goopenmeteo.WindSpeed10MMax,
		goopenmeteo.WindGusts10MMax,
		goopenmeteo.WindDirection10MDominant,
		goopenmeteo.UVIndexMax,
		goopenmeteo.Sunrise,
		goopenmeteo.Sunset,
		goopenmeteo.DaylightDuration,
	}
)

// ForecastDays is how many days are requested and returned, today first.
const ForecastDays = 7

type ForecastResponse struct {
	Current   CurrentForecastResponse  `json:"current"`
	Hourly    []HourlyForecastResponse `json:"hourly"`    // from the current hour to the end of the last day
	Daily     []DailyForecastResponse  `json:"daily"`     // ForecastDays days from today
	Timezone  string                   `json:"timezone"`  // the location's IANA zone, e.g. "Europe/Zagreb"
	Elevation float64                  `json:"elevation"` // metres
	FetchedAt time.Time                `json:"fetchedAt"`
}

type CurrentForecastResponse struct {
	Time                time.Time `json:"time"`
	Temperature         float64   `json:"temperature"`         // °C
	ApparentTemperature float64   `json:"apparentTemperature"` // °C
	Humidity            float64   `json:"humidity"`            // %
	DewPoint            float64   `json:"dewPoint"`            // °C
	Precipitation       float64   `json:"precipitation"`       // mm
	WeatherCode         int       `json:"weatherCode"`         // WMO
	CloudCover          float64   `json:"cloudCover"`          // %
	Pressure            float64   `json:"pressure"`            // hPa at sea level
	WindSpeed           float64   `json:"windSpeed"`           // km/h
	WindDirection       float64   `json:"windDirection"`       // degrees, where it blows from
	WindGusts           float64   `json:"windGusts"`           // km/h
	Visibility          float64   `json:"visibility"`          // m
	UVIndex             float64   `json:"uvIndex"`
	IsDay               bool      `json:"isDay"`
}

type HourlyForecastResponse struct {
	Time                     time.Time `json:"time"`
	Temperature              float64   `json:"temperature"`
	ApparentTemperature      float64   `json:"apparentTemperature"`
	PrecipitationProbability float64   `json:"precipitationProbability"`
	Precipitation            float64   `json:"precipitation"`
	WeatherCode              int       `json:"weatherCode"`
	WindSpeed                float64   `json:"windSpeed"`
	WindDirection            float64   `json:"windDirection"`
	IsDay                    bool      `json:"isDay"`
}

type DailyForecastResponse struct {
	Date                        string    `json:"date"` // YYYY-MM-DD in the location's zone
	TemperatureMin              float64   `json:"temperatureMin"`
	TemperatureMax              float64   `json:"temperatureMax"`
	WeatherCode                 int       `json:"weatherCode"`
	PrecipitationSum            float64   `json:"precipitationSum"`
	PrecipitationProbabilityMax float64   `json:"precipitationProbabilityMax"`
	WindSpeedMax                float64   `json:"windSpeedMax"`
	WindGustsMax                float64   `json:"windGustsMax"`
	WindDirectionDominant       float64   `json:"windDirectionDominant"`
	UVIndexMax                  float64   `json:"uvIndexMax"`
	Sunrise                     time.Time `json:"sunrise"`
	Sunset                      time.Time `json:"sunset"`
	DaylightSeconds             float64   `json:"daylightSeconds"`
}

// NewForecastResponse maps a forecast to the response: hourly from the hour
// containing now, daily from today. Times are placed in the location's own
// zone, so they keep the right offset across a DST change within the week:
// Open-Meteo sends wall-clock times, and goopenmeteo can only give them the
// response's current offset.
func NewForecastResponse(f *goopenmeteo.Forecast, now, fetchedAt time.Time) ForecastResponse {
	zone := zoneOf(f)
	cur := f.Current.Data
	res := ForecastResponse{
		Current: CurrentForecastResponse{
			Time:                zone(f.Current.Time),
			Temperature:         cur[goopenmeteo.Temperature2M],
			ApparentTemperature: cur[goopenmeteo.ApparentTemperature],
			Humidity:            cur[goopenmeteo.RelativeHumidity2M],
			DewPoint:            cur[goopenmeteo.DewPoint2M],
			Precipitation:       cur[goopenmeteo.Precipitation],
			WeatherCode:         int(cur[goopenmeteo.WeatherCode]),
			CloudCover:          cur[goopenmeteo.CloudCover],
			Pressure:            cur[goopenmeteo.PressureMsl],
			WindSpeed:           cur[goopenmeteo.WindSpeed10M],
			WindDirection:       cur[goopenmeteo.WindDirection10M],
			WindGusts:           cur[goopenmeteo.WindGusts10M],
			Visibility:          cur[goopenmeteo.Visibility],
			UVIndex:             cur[goopenmeteo.UVIndex],
			IsDay:               cur[goopenmeteo.IsDay] == 1,
		},
		Hourly:    make([]HourlyForecastResponse, 0, len(f.Hourly.Time)),
		Daily:     make([]DailyForecastResponse, 0, ForecastDays),
		Timezone:  f.Timezone,
		Elevation: f.Elevation,
		FetchedAt: fetchedAt,
	}

	hourly := f.Hourly.Data
	for i, t := range f.Hourly.Time {
		t = zone(t)
		if !t.Add(time.Hour).After(now) { // this hour is already over
			continue
		}
		res.Hourly = append(res.Hourly, HourlyForecastResponse{
			Time:                     t,
			Temperature:              at(hourly[goopenmeteo.Temperature2M], i),
			ApparentTemperature:      at(hourly[goopenmeteo.ApparentTemperature], i),
			PrecipitationProbability: at(hourly[goopenmeteo.PrecipitationProbability], i),
			Precipitation:            at(hourly[goopenmeteo.Precipitation], i),
			WeatherCode:              int(at(hourly[goopenmeteo.WeatherCode], i)),
			WindSpeed:                at(hourly[goopenmeteo.WindSpeed10M], i),
			WindDirection:            at(hourly[goopenmeteo.WindDirection10M], i),
			IsDay:                    at(hourly[goopenmeteo.IsDay], i) >= 1,
		})
	}

	daily, times := f.Daily.Data, f.Daily.Times
	for i, t := range f.Daily.Time {
		if len(res.Daily) == ForecastDays {
			break
		}
		t = zone(t)
		if !t.AddDate(0, 0, 1).After(now) { // this day is already over
			continue
		}
		res.Daily = append(res.Daily, DailyForecastResponse{
			Date:                        t.Format(time.DateOnly),
			TemperatureMin:              at(daily[goopenmeteo.Temperature2MMin], i),
			TemperatureMax:              at(daily[goopenmeteo.Temperature2MMax], i),
			WeatherCode:                 int(at(daily[goopenmeteo.WeatherCode], i)),
			PrecipitationSum:            at(daily[goopenmeteo.PrecipitationSum], i),
			PrecipitationProbabilityMax: at(daily[goopenmeteo.PrecipitationProbabilityMax], i),
			WindSpeedMax:                at(daily[goopenmeteo.WindSpeed10MMax], i),
			WindGustsMax:                at(daily[goopenmeteo.WindGusts10MMax], i),
			WindDirectionDominant:       at(daily[goopenmeteo.WindDirection10MDominant], i),
			UVIndexMax:                  at(daily[goopenmeteo.UVIndexMax], i),
			Sunrise:                     zone(at(times[goopenmeteo.Sunrise], i)),
			Sunset:                      zone(at(times[goopenmeteo.Sunset], i)),
			DaylightSeconds:             at(daily[goopenmeteo.DaylightDuration], i),
		})
	}

	return res
}

// zoneOf re-places wall-clock times in the forecast's IANA zone. With no zone
// (a GMT forecast) or an unknown one, times stay as decoded.
func zoneOf(f *goopenmeteo.Forecast) func(time.Time) time.Time {
	loc, err := time.LoadLocation(f.Timezone)
	if f.Timezone == "" || err != nil {
		return func(t time.Time) time.Time { return t }
	}
	return func(t time.Time) time.Time {
		if t.IsZero() {
			return t
		}
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	}
}

// at returns series[i], or the zero value when the series is shorter than the
// time axis (a variable the response didn't include).
func at[T any](series []T, i int) T {
	if i < len(series) {
		return series[i]
	}
	var zero T
	return zero
}
