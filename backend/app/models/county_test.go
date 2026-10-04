package models

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestCountiesCoverEveryGeoNamesCode(t *testing.T) {
	if len(Counties) != 21 {
		t.Fatalf("%d counties, want 21", len(Counties))
	}
	for i := 1; i <= 21; i++ {
		code := fmt.Sprintf("%02d", i)
		c, ok := Counties[code]
		if !ok || c.Code != code || c.Name == "" || len(c.Areas) < 3 {
			t.Errorf("county %s = %+v", code, c)
		}
	}
}

// The coastal names and NUTS codes are pinned by the recorded Meteoalarm feed
// (testdata/meteoalarm-croatia.json, Task 6).
func TestCountiesUseMeteoalarmsAreaNames(t *testing.T) {
	for code, area := range map[string]string{"03": "Dubrovačko-neretvanska", "08": "Ličko-senjska", "12": "Primorsko-goranska", "13": "Šibensko-kninska", "15": "Splitsko-dalmatinska", "19": "Zadarska"} {
		if !slices.Contains(Counties[code].Areas, area) {
			t.Errorf("county %s areas %v lack %q", code, Counties[code].Areas, area)
		}
	}
	for _, c := range Counties {
		if !slices.ContainsFunc(c.Areas, func(a string) bool { return strings.HasPrefix(a, "HR0") }) {
			t.Errorf("%s has no NUTS3 code", c.Name)
		}
	}
}
