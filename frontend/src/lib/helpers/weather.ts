import type { HourlyForecast } from "$lib/types/forecast";
import { formatHour } from "./format";

export type Precipitation = "none" | "drizzle" | "rain" | "snow" | "storm";

export interface WeatherInfo {
  /** Croatian description of the WMO code. */
  label: string;
  /** A Meteocons icon name (see weather-icon.svelte). */
  icon: string;
  precipitation: Precipitation;
}

type Entry = [label: string, day: string, night: string, precipitation: Precipitation];

// WMO weather interpretation codes, as Open-Meteo reports them.
const codes: Record<number, Entry> = {
  0: ["Vedro", "clear-day", "clear-night", "none"],
  1: ["Pretežno vedro", "partly-cloudy-day", "partly-cloudy-night", "none"],
  2: ["Djelomično oblačno", "partly-cloudy-day", "partly-cloudy-night", "none"],
  3: ["Oblačno", "overcast-day", "overcast-night", "none"],
  45: ["Magla", "fog-day", "fog-night", "none"],
  48: ["Magla s injem", "fog-day", "fog-night", "none"],
  51: ["Slaba rosulja", "partly-cloudy-day-drizzle", "partly-cloudy-night-drizzle", "drizzle"],
  53: ["Rosulja", "drizzle", "drizzle", "drizzle"],
  55: ["Jaka rosulja", "drizzle", "drizzle", "drizzle"],
  56: ["Ledena rosulja", "sleet", "sleet", "drizzle"],
  57: ["Jaka ledena rosulja", "sleet", "sleet", "drizzle"],
  61: ["Slaba kiša", "partly-cloudy-day-rain", "partly-cloudy-night-rain", "rain"],
  63: ["Kiša", "rain", "rain", "rain"],
  65: ["Jaka kiša", "rain", "rain", "rain"],
  66: ["Ledena kiša", "sleet", "sleet", "rain"],
  67: ["Jaka ledena kiša", "sleet", "sleet", "rain"],
  71: ["Slab snijeg", "partly-cloudy-day-snow", "partly-cloudy-night-snow", "snow"],
  73: ["Snijeg", "snow", "snow", "snow"],
  75: ["Jak snijeg", "snow", "snow", "snow"],
  77: ["Zrnati snijeg", "snow", "snow", "snow"],
  80: ["Slabi pljuskovi", "partly-cloudy-day-rain", "partly-cloudy-night-rain", "rain"],
  81: ["Pljuskovi", "rain", "rain", "rain"],
  82: ["Jaki pljuskovi", "rain", "rain", "rain"],
  85: ["Snježni pljuskovi", "partly-cloudy-day-snow", "partly-cloudy-night-snow", "snow"],
  86: ["Jaki snježni pljuskovi", "snow", "snow", "snow"],
  95: ["Grmljavina", "thunderstorms-day-rain", "thunderstorms-night-rain", "storm"],
  96: ["Grmljavina s tučom", "thunderstorms-day-rain", "thunderstorms-night-rain", "storm"],
  99: ["Grmljavina s jakom tučom", "hail", "hail", "storm"],
};

export function weatherInfo(code: number, isDay = true): WeatherInfo {
  const entry = codes[code];
  if (!entry) return { label: "", icon: "not-available", precipitation: "none" };
  const [label, day, night, precipitation] = entry;
  return { label, icon: isDay ? day : night, precipitation };
}

/** An hour counts as wet when rain is likely or measurable. */
const isWet = (h: HourlyForecast) => h.precipitationProbability >= 50 || h.precipitation >= 0.3;

const noun: Record<Precipitation, string> = {
  none: "Kiša",
  drizzle: "Rosulja",
  rain: "Kiša",
  snow: "Snijeg",
  storm: "Grmljavina",
};

const HORIZON = 12;

/** A plain-language line about precipitation over the next 12 hours, starting with the hour that
 *  contains now: "Kiša od 16 h", "Kiša do 15 h", "Bez oborina idućih 12 h". */
export function rainHint(hourly: HourlyForecast[], now: Date, timeZone: string): string {
  const upcoming = hourly.filter((h) => Date.parse(h.time) + 3_600_000 > now.getTime()).slice(0, HORIZON);
  if (upcoming.length === 0) return "";

  const kind = (h: HourlyForecast) => noun[weatherInfo(h.weatherCode).precipitation];

  if (isWet(upcoming[0])) {
    const dry = upcoming.find((h) => !isWet(h));
    return dry ? `${kind(upcoming[0])} do ${formatHour(new Date(dry.time), timeZone)}` : `${kind(upcoming[0])} i dalje`;
  }
  const wet = upcoming.find(isWet);
  return wet ? `${kind(wet)} od ${formatHour(new Date(wet.time), timeZone)}` : `Bez oborina idućih ${HORIZON} h`;
}
