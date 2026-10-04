import type { DailyForecast, Forecast, HourlyForecast } from "$lib/types/forecast";
import { isoDateOf } from "./format";

/** Bumped whenever Forecast changes shape, so an older cached copy is ignored. */
export const CACHE_VERSION = 1;

/** How long a cached forecast is still worth painting while the fresh one loads. */
export const CACHE_MAX_AGE_MS = 12 * 3_600_000;

/** The last forecast shown, kept on the device so the next visit paints at once. */
export interface CachedForecast {
  version: number;
  /** cellKey of the location it is for */
  cell: string;
  savedAt: string;
  forecast: Forecast;
}

export function isCachedForecast(v: unknown): v is CachedForecast {
  if (typeof v !== "object" || v === null) return false;
  const c = v as Record<string, unknown>;
  const f = c.forecast as Record<string, unknown> | null | undefined;
  return (
    c.version === CACHE_VERSION &&
    typeof c.cell === "string" &&
    typeof c.savedAt === "string" &&
    !Number.isNaN(Date.parse(c.savedAt)) &&
    typeof f === "object" &&
    f !== null &&
    typeof f.current === "object" &&
    f.current !== null &&
    Array.isArray(f.hourly) &&
    Array.isArray(f.daily) &&
    typeof f.timezone === "string"
  );
}

/** The cached forecast when it is for cell and recent enough to paint. */
export function usableCache(cache: CachedForecast | undefined, cell: string, now: Date): CachedForecast | undefined {
  if (!cache || cache.cell !== cell) return undefined;
  const age = now.getTime() - Date.parse(cache.savedAt);
  return age >= 0 && age < CACHE_MAX_AGE_MS ? cache : undefined;
}

/** At most count hours, from the one containing now. */
export const upcomingHours = (hourly: HourlyForecast[], now: Date, count: number) =>
  hourly.filter((h) => Date.parse(h.time) + 3_600_000 > now.getTime()).slice(0, count);

/** Today on the location's calendar, or the first day the forecast has. */
export const todayOf = (f: Forecast, now: Date): DailyForecast =>
  f.daily.find((d) => d.date === isoDateOf(now, f.timezone)) ?? f.daily[0];

/** The hours of one day (YYYY-MM-DD on the location's calendar). */
export const hoursOfDay = (f: Forecast, date: string) =>
  f.hourly.filter((h) => isoDateOf(new Date(h.time), f.timezone) === date);

/** The coldest minimum and warmest maximum: the shared scale of the week's range bars. */
export function weekRange(days: DailyForecast[]): [number, number] {
  if (days.length === 0) return [0, 1];
  return [Math.min(...days.map((d) => d.temperatureMin)), Math.max(...days.map((d) => d.temperatureMax))];
}
