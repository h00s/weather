package models

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// WarningLookahead is how far ahead an upcoming warning is already shown: today and tomorrow.
const WarningLookahead = 48 * time.Hour

// MeteoalarmFeed is a country feed from feeds.meteoalarm.org/api/v1/warnings/feeds-<country>.
type MeteoalarmFeed struct {
	Warnings []struct {
		Alert MeteoalarmAlert `json:"alert"`
	} `json:"warnings"`
}

// MeteoalarmAlert is a CAP alert, reduced to the fields the display uses.
type MeteoalarmAlert struct {
	Identifier string           `json:"identifier"`
	MsgType    string           `json:"msgType"`    // Alert, Update or Cancel
	References string           `json:"references"` // "sender,identifier,sent" triplets it replaces
	Info       []MeteoalarmInfo `json:"info"`       // one per language
}

type MeteoalarmInfo struct {
	Language     string    `json:"language"`
	Event        string    `json:"event"`
	Description  string    `json:"description"`
	Onset        time.Time `json:"onset"`
	Expires      time.Time `json:"expires"`
	ResponseType []string  `json:"responseType"`
	Parameter    []struct {
		ValueName string `json:"valueName"`
		Value     string `json:"value"`
	} `json:"parameter"`
	Area []struct {
		AreaDesc string `json:"areaDesc"`
		Geocode  []struct {
			ValueName string `json:"valueName"` // EMMA_ID or NUTS3
			Value     string `json:"value"`
		} `json:"geocode"`
	} `json:"area"`
}

type WarningResponse struct {
	Level       string    `json:"level"` // yellow, orange, red
	Type        string    `json:"type"`  // wind, snow-ice, thunderstorm, fog, high-temperature, …
	Event       string    `json:"event"` // Croatian headline, e.g. "Žuto upozorenje za vjetar"
	Description string    `json:"description"`
	Onset       time.Time `json:"onset"`
	Expires     time.Time `json:"expires"`
}

// Meteoalarm awareness types by number.
var warningTypes = map[int]string{
	1: "wind", 2: "snow-ice", 3: "thunderstorm", 4: "fog", 5: "high-temperature",
	6: "low-temperature", 7: "coastal-event", 8: "forest-fire", 9: "avalanche",
	10: "rain", 12: "flooding", 13: "rain-flood",
}

var warningLevels = map[int]string{2: "yellow", 3: "orange", 4: "red"}

// ActiveWarnings returns the yellow, orange and red warnings for any of areas
// (a county's areaDesc, full name or NUTS3 code; see County.Areas) that are in
// force now or start within WarningLookahead, most severe first. Cancelled
// alerts and alerts a later update replaced are dropped.
func ActiveWarnings(alerts []MeteoalarmAlert, areas []string, now time.Time) []WarningResponse {
	replaced := map[string]bool{}
	for _, a := range alerts {
		for _, ref := range strings.Fields(a.References) {
			if parts := strings.Split(ref, ","); len(parts) == 3 {
				replaced[parts[1]] = true
			}
		}
	}

	type ranked struct {
		WarningResponse
		severity int
	}
	var found []ranked
	seen := map[WarningResponse]bool{}
	for _, a := range alerts {
		if a.MsgType == "Cancel" || replaced[a.Identifier] {
			continue
		}
		info := croatian(a.Info)
		if info == nil || !info.covers(areas) || slices.Contains(info.ResponseType, "AllClear") {
			continue
		}
		if !info.Expires.After(now) || info.Onset.After(now.Add(WarningLookahead)) {
			continue
		}
		severity := info.number("awareness_level")
		level, ok := warningLevels[severity]
		if !ok {
			continue // green: no warning
		}
		w := WarningResponse{
			Level:       level,
			Type:        warningTypes[info.number("awareness_type")],
			Event:       info.Event,
			Description: info.Description,
			Onset:       info.Onset,
			Expires:     info.Expires,
		}
		if seen[w] {
			continue
		}
		seen[w] = true
		found = append(found, ranked{w, severity})
	}

	slices.SortStableFunc(found, func(a, b ranked) int {
		return cmp.Or(cmp.Compare(b.severity, a.severity), a.Onset.Compare(b.Onset))
	})
	out := make([]WarningResponse, len(found))
	for i, f := range found {
		out[i] = f.WarningResponse
	}
	return out
}

// croatian picks the hr-HR info block, else the first.
func croatian(infos []MeteoalarmInfo) *MeteoalarmInfo {
	for i := range infos {
		if infos[i].Language == "hr-HR" {
			return &infos[i]
		}
	}
	if len(infos) > 0 {
		return &infos[0]
	}
	return nil
}

// covers reports whether the info's area is one of areas, by its areaDesc or its
// NUTS3 code. EMMA_IDs are ignored: they share the HR0xx space with NUTS3 codes
// (Splitsko-dalmatinska's EMMA_ID HR027 is Karlovačka's NUTS3), so matching them
// put coastal warnings on inland counties.
func (info *MeteoalarmInfo) covers(areas []string) bool {
	for _, a := range info.Area {
		for _, area := range areas {
			if strings.EqualFold(a.AreaDesc, area) {
				return true
			}
			for _, g := range a.Geocode {
				if g.ValueName == "NUTS3" && strings.EqualFold(g.Value, area) {
					return true
				}
			}
		}
	}
	return false
}

// number reads the leading number of a parameter such as "2; yellow; Moderate".
func (info *MeteoalarmInfo) number(name string) int {
	for _, p := range info.Parameter {
		if p.ValueName == name {
			first, _, _ := strings.Cut(p.Value, ";")
			n, _ := strconv.Atoi(strings.TrimSpace(first))
			return n
		}
	}
	return 0
}
