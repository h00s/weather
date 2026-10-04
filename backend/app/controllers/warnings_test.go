package controllers_test

import (
	"net/http"
	"testing"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

func TestWarningsIndexListsTheCountysWarnings(t *testing.T) {
	app, _ := newApp(t)

	daruvar := raptor.DecodeJSON[[]models.WarningResponse](t, app.TestGet("/api/v1/warnings?lat=45.59&lon=17.22"), http.StatusOK)
	zadar := raptor.DecodeJSON[[]models.WarningResponse](t, app.TestGet("/api/v1/warnings?lat=44.12&lon=15.23"), http.StatusOK)

	if len(daruvar) != 1 || daruvar[0].Level != "yellow" || daruvar[0].Type != "wind" {
		t.Errorf("Daruvar: %+v", daruvar)
	}
	if len(zadar) != 1 || zadar[0].Level != "orange" {
		t.Errorf("Zadar: %+v", zadar)
	}
}

func TestWarningsIndexIsEmptyOutsideCroatia(t *testing.T) {
	app, u := newApp(t)

	got := raptor.DecodeJSON[[]models.WarningResponse](t, app.TestGet("/api/v1/warnings?lat=43.86&lon=18.41"), http.StatusOK)

	if got == nil || len(got) != 0 {
		t.Errorf("Sarajevo: %#v, want []", got)
	}
	if n := u.hits("/api/v1/warnings/feeds-croatia"); n != 0 {
		t.Errorf("Meteoalarm asked %d times for a point outside Croatia", n)
	}
}

func TestWarningsIndexAnswers502WhenMeteoalarmFailsAndNothingIsCached(t *testing.T) {
	app, u := newApp(t)
	u.fail(http.StatusServiceUnavailable)

	raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/warnings?lat=45.59&lon=17.22"), http.StatusBadGateway)
}
