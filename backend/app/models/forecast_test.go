package models

import (
	"testing"
	"time"

	"github.com/h00s/goopenmeteo"
)

// zagreb is how goopenmeteo places a CEST response: a fixed zone at today's offset.
var zagreb = time.FixedZone("GMT+2", 7200)

func hours(start time.Time, n int) []time.Time {
	times := make([]time.Time, n)
	for i := range times {
		times[i] = start.Add(time.Duration(i) * time.Hour)
	}
	return times
}

func seq(n int, from float64) []float64 {
	values := make([]float64, n)
	for i := range values {
		values[i] = from + float64(i)
	}
	return values
}

func testForecast() *goopenmeteo.Forecast {
	midnight := time.Date(2026, 9, 29, 0, 0, 0, 0, zagreb)
	days := []time.Time{midnight, midnight.AddDate(0, 0, 1), midnight.AddDate(0, 0, 2)}
	return &goopenmeteo.Forecast{
		Timezone:  "Europe/Zagreb",
		Elevation: 161,
		Current: goopenmeteo.Current{
			Time: time.Date(2026, 9, 29, 14, 15, 0, 0, zagreb),
			Data: map[string]float64{
				goopenmeteo.Temperature2M:       17.4,
				goopenmeteo.ApparentTemperature: 15.9,
				goopenmeteo.RelativeHumidity2M:  71,
				goopenmeteo.DewPoint2M:          12.1,
				goopenmeteo.Precipitation:       0.2,
				goopenmeteo.WeatherCode:         61,
				goopenmeteo.CloudCover:          90,
				goopenmeteo.PressureMsl:         1013.2,
				goopenmeteo.WindSpeed10M:        12.5,
				goopenmeteo.WindDirection10M:    225,
				goopenmeteo.WindGusts10M:        31,
				goopenmeteo.Visibility:          24000,
				goopenmeteo.UVIndex:             2.6,
				goopenmeteo.IsDay:               1,
			},
		},
		Hourly: goopenmeteo.TimeseriesData{
			Time: hours(midnight, 72),
			Data: map[string][]float64{
				goopenmeteo.Temperature2M:            seq(72, 0),
				goopenmeteo.ApparentTemperature:      seq(72, -2),
				goopenmeteo.PrecipitationProbability: seq(72, 100),
				goopenmeteo.Precipitation:            seq(72, 200),
				goopenmeteo.WeatherCode:              seq(72, 0),
				goopenmeteo.WindSpeed10M:             seq(72, 300),
				goopenmeteo.WindDirection10M:         seq(72, 400),
				goopenmeteo.IsDay:                    seq(72, 0), // 0 at midnight, then 1 and up: "day"
			},
		},
		Daily: goopenmeteo.TimeseriesData{
			Time: days,
			Data: map[string][]float64{
				goopenmeteo.Temperature2MMin:            {9, 8, 7},
				goopenmeteo.Temperature2MMax:            {19, 18, 17},
				goopenmeteo.WeatherCode:                 {61, 3, 0},
				goopenmeteo.PrecipitationSum:            {4.2, 0, 0},
				goopenmeteo.PrecipitationProbabilityMax: {80, 10, 0},
				goopenmeteo.WindSpeed10MMax:             {25, 14, 9},
				goopenmeteo.WindGusts10MMax:             {48, 30, 20},
				goopenmeteo.WindDirection10MDominant:    {200, 45, 90},
				goopenmeteo.UVIndexMax:                  {3.1, 4, 5},
				goopenmeteo.DaylightDuration:            {42600, 42400, 42200},
			},
			Times: map[string][]time.Time{
				goopenmeteo.Sunrise: {days[0].Add(6*time.Hour + 46*time.Minute), days[1].Add(6*time.Hour + 47*time.Minute), days[2].Add(6*time.Hour + 49*time.Minute)},
				goopenmeteo.Sunset:  {days[0].Add(18*time.Hour + 35*time.Minute), days[1].Add(18*time.Hour + 33*time.Minute), days[2].Add(18*time.Hour + 31*time.Minute)},
			},
		},
	}
}

func TestNewForecastResponseMapsCurrentConditions(t *testing.T) {
	now := time.Date(2026, 9, 29, 14, 20, 0, 0, zagreb)
	fetchedAt := now.Add(-time.Minute)

	res := NewForecastResponse(testForecast(), now, fetchedAt)

	want := CurrentForecastResponse{
		Time:                time.Date(2026, 9, 29, 14, 15, 0, 0, zagreb),
		Temperature:         17.4,
		ApparentTemperature: 15.9,
		Humidity:            71,
		DewPoint:            12.1,
		Precipitation:       0.2,
		WeatherCode:         61,
		CloudCover:          90,
		Pressure:            1013.2,
		WindSpeed:           12.5,
		WindDirection:       225,
		WindGusts:           31,
		Visibility:          24000,
		UVIndex:             2.6,
		IsDay:               true,
	}
	got := res.Current
	if !got.Time.Equal(want.Time) {
		t.Errorf("current time = %v, want %v", got.Time, want.Time)
	}
	got.Time = want.Time
	if got != want {
		t.Errorf("current = %+v\nwant      %+v", got, want)
	}
	if res.Timezone != "Europe/Zagreb" || res.Elevation != 161 || !res.FetchedAt.Equal(fetchedAt) {
		t.Errorf("timezone %q, elevation %v, fetchedAt %v", res.Timezone, res.Elevation, res.FetchedAt)
	}
}

func TestNewForecastResponseHourlyStartsAtTheCurrentHour(t *testing.T) {
	now := time.Date(2026, 9, 29, 14, 20, 0, 0, zagreb)

	res := NewForecastResponse(testForecast(), now, now)

	if len(res.Hourly) != 72-14 {
		t.Fatalf("got %d hours, want %d (from 14:00 to the end)", len(res.Hourly), 72-14)
	}
	first := res.Hourly[0]
	if !first.Time.Equal(time.Date(2026, 9, 29, 14, 0, 0, 0, zagreb)) {
		t.Errorf("first hour = %v, want 14:00", first.Time)
	}
	if first.Temperature != 14 || first.ApparentTemperature != 12 || first.PrecipitationProbability != 114 ||
		first.Precipitation != 214 || first.WeatherCode != 14 || first.WindSpeed != 314 || first.WindDirection != 414 {
		t.Errorf("first hour values = %+v, want index 14 of each series", first)
	}
	if !res.Hourly[1].IsDay {
		t.Error("IsDay: a series value of 1 must read as day")
	}
}

func TestNewForecastResponseDailyStartsToday(t *testing.T) {
	now := time.Date(2026, 9, 30, 0, 30, 0, 0, zagreb)

	res := NewForecastResponse(testForecast(), now, now)

	if len(res.Daily) != 2 || res.Daily[0].Date != "2026-09-30" {
		t.Fatalf("daily = %+v, want 2 days from 2026-09-30", res.Daily)
	}
	d := res.Daily[0]
	if d.TemperatureMin != 8 || d.TemperatureMax != 18 || d.WeatherCode != 3 || d.PrecipitationProbabilityMax != 10 ||
		d.WindSpeedMax != 14 || d.WindGustsMax != 30 || d.WindDirectionDominant != 45 || d.UVIndexMax != 4 || d.DaylightSeconds != 42400 {
		t.Errorf("day = %+v", d)
	}
	if !d.Sunrise.Equal(time.Date(2026, 9, 30, 6, 47, 0, 0, zagreb)) || !d.Sunset.Equal(time.Date(2026, 9, 30, 18, 33, 0, 0, zagreb)) {
		t.Errorf("sun = %v – %v", d.Sunrise, d.Sunset)
	}
}

func TestNewForecastResponseToleratesShortSeries(t *testing.T) {
	f := testForecast()
	f.Hourly.Data[goopenmeteo.Temperature2M] = []float64{1, 2}
	delete(f.Daily.Times, goopenmeteo.Sunset)
	now := time.Date(2026, 9, 29, 0, 10, 0, 0, zagreb)

	res := NewForecastResponse(f, now, now)

	if res.Hourly[1].Temperature != 2 || res.Hourly[5].Temperature != 0 {
		t.Errorf("a short series should read as zeros past its end: %v, %v", res.Hourly[1].Temperature, res.Hourly[5].Temperature)
	}
	if !res.Daily[0].Sunset.IsZero() {
		t.Errorf("missing sunset = %v, want the zero time", res.Daily[0].Sunset)
	}
}

// Open-Meteo labels the whole series with the offset in force when it was asked
// (utc_offset_seconds), straight through a DST change: after 25 Oct, "05:00" at
// GMT+2 is 03:00 UTC, which the CET clock reads as 04:00. goopenmeteo's fixed zone
// already makes these the right instants; re-placing the labels in Europe/Zagreb
// would shift every hour after the switch by one.
func TestNewForecastResponseKeepsUpstreamInstantsAcrossDST(t *testing.T) {
	start := time.Date(2026, 10, 24, 0, 0, 0, 0, zagreb)
	f := &goopenmeteo.Forecast{
		Timezone: "Europe/Zagreb",
		Hourly: goopenmeteo.TimeseriesData{
			Time: hours(start, 48), // labelled 24 Oct 00:00 … 25 Oct 23:00, all at GMT+2
			Data: map[string][]float64{goopenmeteo.Temperature2M: seq(48, 0)},
		},
	}
	now := time.Date(2026, 10, 24, 0, 30, 0, 0, zagreb)

	res := NewForecastResponse(f, now, now)

	want := time.Date(2026, 10, 25, 3, 0, 0, 0, time.UTC) // the hour labelled "2026-10-25T05:00"
	if got := res.Hourly[29].Time; !got.Equal(want) {
		t.Errorf("hour labelled 25 Oct 05:00 = %v, want %v", got.UTC(), want)
	}
}

// Every hour is a distinct, later instant, even across the switch to CEST, when
// 28 Mar 2027 02:00 does not exist on the local clock: the frontend keys its
// lists on the time, and a repeated one stops the page from rendering.
func TestNewForecastResponseHoursStrictlyIncreaseAcrossDST(t *testing.T) {
	cet := time.FixedZone("GMT+1", 3600)
	start := time.Date(2027, 3, 27, 0, 0, 0, 0, cet)
	f := &goopenmeteo.Forecast{
		Timezone: "Europe/Zagreb",
		Hourly: goopenmeteo.TimeseriesData{
			Time: hours(start, 48), // labelled 27 Mar 00:00 … 28 Mar 23:00, all at GMT+1
			Data: map[string][]float64{goopenmeteo.Temperature2M: seq(48, 0)},
		},
	}

	res := NewForecastResponse(f, start, start)

	if len(res.Hourly) != 48 {
		t.Fatalf("got %d hours, want 48", len(res.Hourly))
	}
	for i := 1; i < len(res.Hourly); i++ {
		if !res.Hourly[i].Time.After(res.Hourly[i-1].Time) {
			t.Fatalf("hour %d (%v) is not after hour %d (%v)", i, res.Hourly[i].Time.UTC(), i-1, res.Hourly[i-1].Time.UTC())
		}
	}
}
