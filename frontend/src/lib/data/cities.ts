import type { SavedLocation } from "$lib/types/location";

const city = (name: string, county: string, latitude: number, longitude: number): SavedLocation => ({
  kind: "place",
  name,
  county,
  latitude,
  longitude,
});

/** One city per county (Zagreb stands for both Grad Zagreb and Zagrebačka), largest first: the
 *  choices offered before any search. Coordinates are GeoNames', rounded like every location. */
export const cities: SavedLocation[] = [
  city("Zagreb", "Grad Zagreb", 45.81, 15.98),
  city("Split", "Splitsko-dalmatinska županija", 43.51, 16.44),
  city("Rijeka", "Primorsko-goranska županija", 45.33, 14.44),
  city("Osijek", "Osječko-baranjska županija", 45.55, 18.69),
  city("Zadar", "Zadarska županija", 44.12, 15.23),
  city("Pula", "Istarska županija", 44.87, 13.85),
  city("Slavonski Brod", "Brodsko-posavska županija", 45.16, 18.02),
  city("Karlovac", "Karlovačka županija", 45.49, 15.55),
  city("Varaždin", "Varaždinska županija", 46.3, 16.34),
  city("Šibenik", "Šibensko-kninska županija", 43.73, 15.89),
  city("Sisak", "Sisačko-moslavačka županija", 45.47, 16.38),
  city("Vinkovci", "Vukovarsko-srijemska županija", 45.29, 18.8),
  city("Dubrovnik", "Dubrovačko-neretvanska županija", 42.64, 18.11),
  city("Bjelovar", "Bjelovarsko-bilogorska županija", 45.9, 16.85),
  city("Koprivnica", "Koprivničko-križevačka županija", 46.16, 16.83),
  city("Požega", "Požeško-slavonska županija", 45.34, 17.69),
  city("Čakovec", "Međimurska županija", 46.38, 16.43),
  city("Virovitica", "Virovitičko-podravska županija", 45.83, 17.39),
  city("Gospić", "Ličko-senjska županija", 44.55, 15.37),
  city("Krapina", "Krapinsko-zagorska županija", 46.16, 15.87),
];
