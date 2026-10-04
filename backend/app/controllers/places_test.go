package controllers_test

import (
	"net/http"
	"strings"
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

func TestPlacesIndexSearchesCroatianSettlements(t *testing.T) {
	app, u := newApp(t)

	got := raptor.DecodeJSON[[]models.PlaceResponse](t, app.TestGet("/api/v1/places?q=%20Daru%20"), http.StatusOK)

	if len(got) != 2 || got[0].Name != "Daruvar" || got[1].Name != "Daruvarski Vinogradi" {
		t.Errorf("got %+v", got)
	}
	q := u.query("/search")
	if q.Get("name") != "Daru" || q.Get("language") != "hr" || q.Get("countryCode") != "HR" {
		t.Errorf("geocoding query = %v", q)
	}
}

func TestPlacesIndexCachesAQueryWhateverItsCase(t *testing.T) {
	app, u := newApp(t)

	for _, q := range []string{"Daru", "daru", "DARU"} {
		raptor.DecodeJSON[[]models.PlaceResponse](t, app.TestGet("/api/v1/places?q="+q), http.StatusOK)
	}

	if n := u.hits("/search"); n != 1 {
		t.Errorf("geocoding asked %d times, want 1", n)
	}
}

func TestPlacesIndexRejectsTooShortOrTooLongQueries(t *testing.T) {
	app, u := newApp(t)

	for _, q := range []string{"", "q=", "q=a", "q=%20%20a%20", "q=" + strings.Repeat("a", 61)} {
		raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/places?"+q), http.StatusBadRequest)
	}
	if n := u.hits("/search"); n != 0 {
		t.Errorf("geocoding asked %d times for invalid queries", n)
	}
}

func TestPlacesIndexAnswers502WhenGeocodingFails(t *testing.T) {
	app, u := newApp(t)
	u.fail(http.StatusInternalServerError)

	raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/places?q=Daru"), http.StatusBadGateway)
}
