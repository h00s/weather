import { describe, expect, it } from "vitest";
import type { DailyForecast, Forecast, HourlyForecast } from "$lib/types/forecast";
import {
  CACHE_MAX_AGE_MS,
  CACHE_VERSION,
  hoursOfDay,
  isCachedForecast,
  todayOf,
  upcomingHours,
  usableCache,
  weekRange,
  type CachedForecast,
} from "./forecast";

const hour = (time: string): HourlyForecast => ({
  time,
  temperature: 10,
  apparentTemperature: 9,
  precipitationProbability: 0,
  precipitation: 0,
  weatherCode: 0,
  windSpeed: 5,
  windDirection: 90,
  isDay: false,
});

const day = (date: string, temperatureMin: number, temperatureMax: number): DailyForecast => ({
  date,
  temperatureMin,
  temperatureMax,
  weatherCode: 0,
  precipitationSum: 0,
  precipitationProbabilityMax: 0,
  windSpeedMax: 10,
  windGustsMax: 20,
  windDirectionDominant: 90,
  uvIndexMax: 3,
  sunrise: `${date}T06:46:00+02:00`,
  sunset: `${date}T18:35:00+02:00`,
  daylightSeconds: 42_600,
});

const forecast = (): Forecast => ({
  current: {
    time: "2026-09-29T23:15:00+02:00",
    temperature: 12,
    apparentTemperature: 11,
    humidity: 70,
    dewPoint: 7,
    precipitation: 0,
    weatherCode: 0,
    cloudCover: 0,
    pressure: 1013,
    windSpeed: 5,
    windDirection: 90,
    windGusts: 9,
    visibility: 20_000,
    uvIndex: 0,
    isDay: false,
  },
  hourly: ["2026-09-29T22:00:00+02:00", "2026-09-29T23:00:00+02:00", "2026-09-30T00:00:00+02:00", "2026-09-30T01:00:00+02:00"].map(hour),
  daily: [day("2026-09-29", 9, 19), day("2026-09-30", 7, 21)],
  timezone: "Europe/Zagreb",
  elevation: 161,
  fetchedAt: "2026-09-29T23:10:00+02:00",
});

describe("upcomingHours", () => {
  it("starts at the hour containing now", () => {
    const got = upcomingHours(forecast().hourly, new Date("2026-09-29T23:20:00+02:00"), 2);
    expect(got.map((h) => h.time)).toEqual(["2026-09-29T23:00:00+02:00", "2026-09-30T00:00:00+02:00"]);
  });
});

describe("todayOf", () => {
  it("picks today on the location's calendar", () => {
    expect(todayOf(forecast(), new Date("2026-09-29T23:30:00+02:00")).date).toBe("2026-09-29");
    // 22:30 UTC is already the 30th in Zagreb, whatever the device's zone.
    expect(todayOf(forecast(), new Date("2026-09-29T22:30:00Z")).date).toBe("2026-09-30");
  });

  it("falls back to the first day", () => {
    expect(todayOf(forecast(), new Date("2026-10-05T12:00:00+02:00")).date).toBe("2026-09-29");
  });
});

describe("hoursOfDay", () => {
  it("keeps one calendar day's hours", () => {
    expect(hoursOfDay(forecast(), "2026-09-30").map((h) => h.time)).toEqual(["2026-09-30T00:00:00+02:00", "2026-09-30T01:00:00+02:00"]);
  });
});

describe("weekRange", () => {
  it("spans the coldest minimum to the warmest maximum", () => {
    expect(weekRange(forecast().daily)).toEqual([7, 21]);
    expect(weekRange([])).toEqual([0, 1]);
  });
});

describe("the cached forecast", () => {
  const cache = (savedAt: string, cell = "45.59,17.23"): CachedForecast => ({ version: CACHE_VERSION, cell, savedAt, forecast: forecast() });
  const now = new Date("2026-09-29T23:30:00+02:00");

  it("is painted only for the same cell while it is recent", () => {
    const fresh = cache("2026-09-29T22:30:00+02:00");
    expect(usableCache(fresh, "45.59,17.23", now)).toBe(fresh);
    expect(usableCache(fresh, "43.51,16.44", now)).toBeUndefined();
    expect(usableCache(cache(new Date(now.getTime() - CACHE_MAX_AGE_MS - 1).toISOString()), "45.59,17.23", now)).toBeUndefined();
    expect(usableCache(cache("2026-09-30T05:00:00+02:00"), "45.59,17.23", now)).toBeUndefined(); // from the future
    expect(usableCache(undefined, "45.59,17.23", now)).toBeUndefined();
  });

  it("is recognized only in this version's shape", () => {
    expect(isCachedForecast(cache("2026-09-29T22:30:00+02:00"))).toBe(true);
    expect(isCachedForecast({ ...cache("2026-09-29T22:30:00+02:00"), version: 0 })).toBe(false);
    expect(isCachedForecast({ ...cache("not a date") })).toBe(false);
    expect(isCachedForecast({ ...cache("2026-09-29T22:30:00+02:00"), forecast: { current: {} } })).toBe(false);
    expect(isCachedForecast(null)).toBe(false);
  });
});
