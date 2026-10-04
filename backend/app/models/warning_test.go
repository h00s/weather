package models

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

var warnNow = time.Date(2026, 9, 29, 12, 0, 0, 0, zagreb)

// alert builds a Meteoalarm alert with Croatian and English info blocks.
func alert(id, msgType, area, level, kind string, onset, expires time.Time, references string) MeteoalarmAlert {
	info := func(lang, event string) string {
		return `{"language":"` + lang + `","event":"` + event + `","description":"Opis","onset":"` + onset.Format(time.RFC3339) +
			`","expires":"` + expires.Format(time.RFC3339) + `","responseType":["Monitor"],
			"parameter":[{"valueName":"awareness_level","value":"` + level + `"},{"valueName":"awareness_type","value":"` + kind + `"}],
			"area":[{"areaDesc":"` + area + `","geocode":[{"valueName":"EMMA_ID","value":"HR005"},{"valueName":"NUTS3","value":"HR021"}]}]}`
	}
	var a MeteoalarmAlert
	raw := `{"identifier":"` + id + `","msgType":"` + msgType + `","references":"` + references + `","info":[` +
		info("en-GB", "Moderate Wind Warning") + `,` + info("hr-HR", "Žuto upozorenje za vjetar") + `]}`
	if err := json.Unmarshal([]byte(raw), &a); err != nil {
		panic(err)
	}
	return a
}

func TestActiveWarningsKeepsCurrentAndUpcomingWarningsForTheArea(t *testing.T) {
	alerts := []MeteoalarmAlert{
		alert("a", "Alert", "Bjelovarsko-bilogorska", "2; yellow; Moderate", "1; Wind", warnNow.Add(-time.Hour), warnNow.Add(6*time.Hour), ""),
		alert("b", "Alert", "Bjelovarsko-bilogorska", "3; orange; Severe", "3; Thunderstorm", warnNow.Add(3*time.Hour), warnNow.Add(9*time.Hour), ""),
	}

	got := ActiveWarnings(alerts, []string{"Bjelovarsko-bilogorska"}, warnNow)

	if len(got) != 2 {
		t.Fatalf("got %d warnings, want 2: %+v", len(got), got)
	}
	// The most severe first.
	if got[0].Level != "orange" || got[0].Type != "thunderstorm" || got[1].Level != "yellow" || got[1].Type != "wind" {
		t.Errorf("warnings = %+v", got)
	}
	if got[1].Event != "Žuto upozorenje za vjetar" || got[1].Description != "Opis" {
		t.Errorf("want the Croatian text, got %+v", got[1])
	}
	if !got[1].Onset.Equal(warnNow.Add(-time.Hour)) || !got[1].Expires.Equal(warnNow.Add(6*time.Hour)) {
		t.Errorf("times = %v..%v", got[1].Onset, got[1].Expires)
	}
}

func TestActiveWarningsMatchesAreaByNameOrGeocode(t *testing.T) {
	alerts := []MeteoalarmAlert{alert("a", "Alert", "Bjelovarsko-bilogorska", "2; yellow; Moderate", "1; Wind", warnNow, warnNow.Add(time.Hour), "")}

	for _, area := range []string{"bjelovarsko-bilogorska", "HR021", "HR005"} {
		if got := ActiveWarnings(alerts, []string{area}, warnNow); len(got) != 1 {
			t.Errorf("area %q matched %d warnings, want 1", area, len(got))
		}
	}
	if got := ActiveWarnings(alerts, []string{"Zadarska"}, warnNow); len(got) != 0 {
		t.Errorf("another area matched %d warnings, want 0", len(got))
	}
}

func TestActiveWarningsDropsWhatIsNotAWarning(t *testing.T) {
	area := "Bjelovarsko-bilogorska"
	alerts := []MeteoalarmAlert{
		alert("green", "Alert", area, "1; green; Minor", "1; Wind", warnNow, warnNow.Add(time.Hour), ""),
		alert("expired", "Alert", area, "2; yellow; Moderate", "1; Wind", warnNow.Add(-3*time.Hour), warnNow.Add(-time.Hour), ""),
		alert("far", "Alert", area, "2; yellow; Moderate", "1; Wind", warnNow.Add(50*time.Hour), warnNow.Add(60*time.Hour), ""),
		alert("cancel", "Cancel", area, "2; yellow; Moderate", "1; Wind", warnNow, warnNow.Add(time.Hour), ""),
		alert("old", "Alert", area, "3; orange; Severe", "10; Rain", warnNow, warnNow.Add(time.Hour), ""),
		alert("new", "Update", area, "2; yellow; Moderate", "10; Rain", warnNow, warnNow.Add(time.Hour), "https://meteo.hr,old,2026-09-29T10:00:00+02:00"),
	}

	got := ActiveWarnings(alerts, []string{area}, warnNow)

	if len(got) != 1 || got[0].Level != "yellow" || got[0].Type != "rain" {
		t.Errorf("want only the update that replaced 'old', got %+v", got)
	}
}

func TestActiveWarningsCollapsesDuplicates(t *testing.T) {
	area := "Bjelovarsko-bilogorska"
	alerts := []MeteoalarmAlert{
		alert("a", "Alert", area, "2; yellow; Moderate", "1; Wind", warnNow, warnNow.Add(time.Hour), ""),
		alert("b", "Alert", area, "2; yellow; Moderate", "1; Wind", warnNow, warnNow.Add(time.Hour), ""),
	}

	if got := ActiveWarnings(alerts, []string{area}, warnNow); len(got) != 1 {
		t.Errorf("got %d warnings, want 1", len(got))
	}
}

// A real Croatian feed (2026-09-28), to pin the structs to Meteoalarm's format.
func TestActiveWarningsOnARecordedFeed(t *testing.T) {
	raw, err := os.ReadFile("testdata/meteoalarm-croatia.json")
	if err != nil {
		t.Fatal(err)
	}
	var feed MeteoalarmFeed
	if err := json.Unmarshal(raw, &feed); err != nil {
		t.Fatal(err)
	}
	alerts := make([]MeteoalarmAlert, len(feed.Warnings))
	for i, w := range feed.Warnings {
		alerts[i] = w.Alert
	}
	now := time.Date(2026, 9, 26, 8, 0, 0, 0, zagreb)

	got := ActiveWarnings(alerts, []string{"Zadarska"}, now)

	if len(got) != 1 {
		t.Fatalf("got %d warnings for Zadarska, want 1: %+v", len(got), got)
	}
	want := WarningResponse{
		Level:   "yellow",
		Type:    "wind",
		Event:   "Žuto upozorenje za vjetar",
		Onset:   time.Date(2026, 9, 26, 6, 35, 57, 0, zagreb),
		Expires: time.Date(2026, 9, 26, 10, 59, 59, 0, zagreb),
	}
	w := got[0]
	if w.Level != want.Level || w.Type != want.Type || w.Event != want.Event || !w.Onset.Equal(want.Onset) || !w.Expires.Equal(want.Expires) {
		t.Errorf("warning = %+v, want %+v", w, want)
	}
	if len(ActiveWarnings(alerts, []string{"Bjelovarsko-bilogorska"}, now)) != 0 {
		t.Error("the recorded feed has no warnings for Bjelovarsko-bilogorska")
	}
}

// The app shows today and tomorrow: a warning starting 30 hours ahead is listed.
func TestActiveWarningsKeepsTomorrowsWarnings(t *testing.T) {
	area := "Bjelovarsko-bilogorska"
	alerts := []MeteoalarmAlert{alert("tomorrow", "Alert", area, "2; yellow; Moderate", "1; Wind", warnNow.Add(30*time.Hour), warnNow.Add(36*time.Hour), "")}

	if got := ActiveWarnings(alerts, []string{area}, warnNow); len(got) != 1 {
		t.Errorf("got %d warnings, want tomorrow's", len(got))
	}
}

func TestActiveWarningsMatchesAnyOfACountysAreas(t *testing.T) {
	alerts := []MeteoalarmAlert{alert("a", "Alert", "Bjelovarsko-bilogorska", "2; yellow; Moderate", "1; Wind", warnNow, warnNow.Add(time.Hour), "")}

	if got := ActiveWarnings(alerts, Counties["01"].Areas, warnNow); len(got) != 1 {
		t.Errorf("Bjelovarsko-bilogorska matched %d warnings, want 1", len(got))
	}
	if got := ActiveWarnings(alerts, Counties["19"].Areas, warnNow); len(got) != 0 {
		t.Errorf("Zadarska matched %d warnings, want 0", len(got))
	}
}
