package controllers_test

import (
	"net/http"
	"testing"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

func TestPlacesNearestNamesTheSettlement(t *testing.T) {
	app, _ := newApp(t)

	got := raptor.DecodeJSON[models.PlaceResponse](t, app.TestGet("/api/v1/places/nearest?lat=45.59&lon=17.22"), http.StatusOK)

	if got.Name != "Daruvar" || got.County != "Bjelovarsko-bilogorska županija" {
		t.Errorf("got %+v", got)
	}
}

func TestPlacesNearestIs404AbroadAnd400ForBadCoordinates(t *testing.T) {
	app, _ := newApp(t)

	raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/places/nearest?lat=43.86&lon=18.41"), http.StatusNotFound)
	raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/places/nearest?lat=x&lon=18"), http.StatusBadRequest)
}
