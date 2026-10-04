import type { AirQuality, AQILevel, PollenLevel, PollenType } from "$lib/types/airquality";

/** The European Air Quality Index bands, as the EEA names them in Croatian. */
export const aqiLabels: Record<AQILevel, string> = {
  good: "Dobra",
  fair: "Prihvatljiva",
  moderate: "Umjerena",
  poor: "Loša",
  very_poor: "Vrlo loša",
  extremely_poor: "Izrazito loša",
};

export const pollenNames: Record<PollenType, string> = {
  alder: "Joha",
  birch: "Breza",
  olive: "Maslina",
  grass: "Trave",
  mugwort: "Pelin",
  ragweed: "Ambrozija",
};

export const pollenLevelLabels: Record<PollenLevel, string> = {
  none: "nema",
  low: "nisko",
  moderate: "umjereno",
  high: "visoko",
  very_high: "vrlo visoko",
};

const severity: Record<PollenLevel, number> = { none: 0, low: 1, moderate: 2, high: 3, very_high: 4 };

/** Today's pollen in the air, the highest first; empty outside the season. */
export const activePollen = (aq: AirQuality) =>
  aq.pollen.filter((p) => p.level !== "none").sort((a, b) => severity[b.level] - severity[a.level]);
