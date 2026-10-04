package models

import (
	"encoding/json"
	"slices"
	"testing"
)

// Recorded from geocoding-api.open-meteo.com (?name=Daruvar&language=hr&countryCode=HR),
// plus a Bosnian and an abandoned place the filter must drop.
const daruvarSearch = `{"results":[
 {"id":3202184,"name":"Daruvar","latitude":45.59056,"longitude":17.225,"feature_code":"PPLA2","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija","population":7440},
 {"id":11500092,"name":"Daruvar","latitude":45.58507,"longitude":17.2114,"feature_code":"AIRF","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija"},
 {"id":12509853,"name":"Daruvarski Vinogradi","latitude":45.60251,"longitude":17.25084,"feature_code":"PPL","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija","population":182},
 {"id":1,"name":"Daruvarac","latitude":44.5,"longitude":18.1,"feature_code":"PPL","country_code":"BA","admin1":"Federacija Bosne i Hercegovine"},
 {"id":2,"name":"Daruvarska Stara","latitude":45.6,"longitude":17.2,"feature_code":"PPLQ","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija"}
],"generationtime_ms":0.31}`

func TestSearchResultsKeepsCroatianSettlementsInOrder(t *testing.T) {
	var g GeocodingResponse
	if err := json.Unmarshal([]byte(daruvarSearch), &g); err != nil {
		t.Fatal(err)
	}

	want := []PlaceResponse{
		{Name: "Daruvar", County: "Bjelovarsko-bilogorska županija", Latitude: 45.59056, Longitude: 17.225},
		{Name: "Daruvarski Vinogradi", County: "Bjelovarsko-bilogorska županija", Latitude: 45.60251, Longitude: 17.25084},
	}
	if got := SearchResults(g, 10); !slices.Equal(got, want) {
		t.Errorf("got %+v\nwant %+v", got, want)
	}
	if got := SearchResults(g, 1); len(got) != 1 || got[0].Name != "Daruvar" {
		t.Errorf("limit 1: %+v", got)
	}
	if got := SearchResults(GeocodingResponse{}, 10); got == nil || len(got) != 0 {
		t.Errorf("no results = %#v, want an empty, non-nil slice", got)
	}
}
