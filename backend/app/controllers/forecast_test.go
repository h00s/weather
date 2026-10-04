package controllers_test

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

func TestForecastShowReturnsTheForecastForTheCell(t *testing.T) {
	app, u := newApp(t)

	res := raptor.DecodeJSON[models.ForecastResponse](t, app.TestGet("/api/v1/forecast?lat=45.5936&lon=17.2251"), http.StatusOK)

	c := res.Current
	if c.Temperature != 17.4 || c.WeatherCode != 61 || !c.IsDay || c.WindDirection != 225 {
		t.Errorf("current = %+v", c)
	}
	if len(res.Daily) != models.ForecastDays || len(res.Hourly) == 0 || res.Timezone != "Europe/Zagreb" || res.FetchedAt.IsZero() {
		t.Errorf("daily %d, hourly %d, timezone %q, fetchedAt %v", len(res.Daily), len(res.Hourly), res.Timezone, res.FetchedAt)
	}
	if first, now := res.Hourly[0].Time, time.Now(); first.After(now) || !first.Add(time.Hour).After(now) {
		t.Errorf("hourly starts at %v, want the hour containing %v", first, now)
	}

	q := u.query("/forecast")
	lat, _ := strconv.ParseFloat(q.Get("latitude"), 64)
	lon, _ := strconv.ParseFloat(q.Get("longitude"), 64)
	if lat != 45.59 || lon != 17.23 {
		t.Errorf("Open-Meteo got %v,%v; want the rounded cell 45.59,17.23", lat, lon)
	}
	if q.Get("timezone") != "auto" || q.Get("forecast_days") != "7" || q.Get("daily") == "" || q.Get("current") == "" {
		t.Errorf("Open-Meteo query = %v", q)
	}
}

func TestForecastShowSharesOneUpstreamRequestPerCell(t *testing.T) {
	app, u := newApp(t)

	for _, q := range []string{"lat=45.5936&lon=17.2251", "lat=45.5912&lon=17.2263"} {
		raptor.DecodeJSON[models.ForecastResponse](t, app.TestGet("/api/v1/forecast?"+q), http.StatusOK)
	}

	if n := u.hits("/forecast"); n != 1 {
		t.Errorf("Open-Meteo was asked %d times for one cell, want 1", n)
	}
}

func TestForecastShowRejectsBadCoordinates(t *testing.T) {
	app, u := newApp(t)

	for _, q := range []string{"", "lat=45.59", "lat=abc&lon=17", "lat=91&lon=17", "lat=45&lon=-180.5", "lat=NaN&lon=17"} {
		body := raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/forecast?"+q), http.StatusBadRequest)
		if body.Code != http.StatusBadRequest || body.Message == "" {
			t.Errorf("%q: body = %+v", q, body)
		}
	}
	if n := u.hits("/forecast"); n != 0 {
		t.Errorf("Open-Meteo was asked %d times for invalid coordinates", n)
	}
}

func TestForecastShowAnswers502WhenOpenMeteoFailsAndNothingIsCached(t *testing.T) {
	app, u := newApp(t)
	u.fail(http.StatusServiceUnavailable)

	body := raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/forecast?lat=45.59&lon=17.23"), http.StatusBadGateway)

	if body.Code != http.StatusBadGateway {
		t.Errorf("body = %+v", body)
	}
}
