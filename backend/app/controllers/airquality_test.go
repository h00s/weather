package controllers_test

import (
	"net/http"
	"testing"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

func TestAirQualityShowReturnsAQIAndPollen(t *testing.T) {
	app, u := newApp(t)

	res := raptor.DecodeJSON[models.AirQualityResponse](t, app.TestGet("/api/v1/air-quality?lat=45.5936&lon=17.2251"), http.StatusOK)

	if res.AQI.Value != 43 || res.AQI.Level != "moderate" {
		t.Errorf("aqi = %+v", res.AQI)
	}
	var ragweed models.PollenResponse
	for _, p := range res.Pollen {
		if p.Type == "ragweed" {
			ragweed = p
		}
	}
	if ragweed.TodayMax != 23.5 || ragweed.Level != "high" || len(res.Pollen) != len(models.PollenTypes) {
		t.Errorf("pollen = %+v", res.Pollen)
	}
	q := u.query("/air-quality")
	if q.Get("timezone") != "auto" || q.Get("forecast_days") != "1" || q.Get("latitude") == "" {
		t.Errorf("air quality query = %v", q)
	}
}

func TestAirQualityShowAnswers502WhenOpenMeteoFailsAndNothingIsCached(t *testing.T) {
	app, u := newApp(t)
	u.fail(http.StatusServiceUnavailable)

	raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/air-quality?lat=45.59&lon=17.22"), http.StatusBadGateway)
}
