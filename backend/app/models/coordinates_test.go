package models

import "testing"

func TestParseCoordinates(t *testing.T) {
	ok := []struct {
		lat, lon string
		want     Coordinates
	}{
		{"45.5936", "17.2251", Coordinates{45.5936, 17.2251}},
		{"-33.9", "151.2", Coordinates{-33.9, 151.2}},
		{"90", "-180", Coordinates{90, -180}},
	}
	for _, c := range ok {
		got, err := ParseCoordinates(c.lat, c.lon)
		if err != nil || got != c.want {
			t.Errorf("ParseCoordinates(%q, %q) = %v, %v; want %v", c.lat, c.lon, got, err, c.want)
		}
	}

	bad := []struct{ lat, lon string }{
		{"", "17"}, {"45", ""}, {"abc", "17"}, {"45", "1e999"}, {"NaN", "17"},
		{"Inf", "17"}, {"91", "17"}, {"45", "180.5"}, {"-90.01", "0"},
	}
	for _, c := range bad {
		if got, err := ParseCoordinates(c.lat, c.lon); err == nil {
			t.Errorf("ParseCoordinates(%q, %q) = %v, want an error", c.lat, c.lon, got)
		}
	}
}

func TestRoundedSnapsToAKilometreCell(t *testing.T) {
	got := Coordinates{45.5936, 17.2251}.Rounded()
	if got != (Coordinates{45.59, 17.23}) {
		t.Errorf("Rounded = %v, want {45.59 17.23}", got)
	}
	if k := got.Key(); k != "45.59,17.23" {
		t.Errorf("Key = %q", k)
	}
	if a, b := (Coordinates{45.5912, 17.2263}).Key(), (Coordinates{45.5936, 17.2251}).Key(); a != b {
		t.Errorf("neighbours in one cell got keys %q and %q", a, b)
	}
	if k := (Coordinates{-0.001, 0.004}).Key(); k != "0.00,0.00" {
		t.Errorf("Key = %q, want no negative zero", k)
	}
}
