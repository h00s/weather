package models

// County is one of Croatia's 21 counties (županije). Areas name it as
// Meteoalarm's warning areas may: the areaDesc DHMZ uses ("Zadarska"), the full
// name, and its NUTS 2021 code. ActiveWarnings matches any of them.
type County struct {
	Code  string // GeoNames admin1 code
	Name  string // e.g. "Bjelovarsko-bilogorska županija"
	Areas []string
}

func county(code, short, nuts string) County {
	name := short + " županija"
	return County{Code: code, Name: name, Areas: []string{short, name, nuts}}
}

// Counties by GeoNames admin1 code. Coastal areaDesc and NUTS codes are checked
// against a recorded feed; inland ones follow the same pattern. If a live
// warning for an inland county ever fails to show, compare its areaDesc here.
var Counties = map[string]County{
	"01": county("01", "Bjelovarsko-bilogorska", "HR021"),
	"02": county("02", "Brodsko-posavska", "HR024"),
	"03": county("03", "Dubrovačko-neretvanska", "HR037"),
	"04": county("04", "Istarska", "HR036"),
	"05": county("05", "Karlovačka", "HR027"),
	"06": county("06", "Koprivničko-križevačka", "HR063"),
	"07": county("07", "Krapinsko-zagorska", "HR064"),
	"08": county("08", "Ličko-senjska", "HR032"),
	"09": county("09", "Međimurska", "HR061"),
	"10": county("10", "Osječko-baranjska", "HR025"),
	"11": county("11", "Požeško-slavonska", "HR023"),
	"12": county("12", "Primorsko-goranska", "HR031"),
	"13": county("13", "Šibensko-kninska", "HR034"),
	"14": county("14", "Sisačko-moslavačka", "HR028"),
	"15": county("15", "Splitsko-dalmatinska", "HR035"),
	"16": county("16", "Varaždinska", "HR062"),
	"17": county("17", "Virovitičko-podravska", "HR022"),
	"18": county("18", "Vukovarsko-srijemska", "HR026"),
	"19": county("19", "Zadarska", "HR033"),
	"20": county("20", "Zagrebačka", "HR065"),
	"21": {Code: "21", Name: "Grad Zagreb", Areas: []string{"Grad Zagreb", "Zagreb", "HR050"}},
}
