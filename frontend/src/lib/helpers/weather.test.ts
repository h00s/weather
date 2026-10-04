import { describe, expect, it } from "vitest";
import type { HourlyForecast } from "$lib/types/forecast";
import { rainHint, weatherInfo } from "./weather";

const tz = "Europe/Zagreb";

describe("weatherInfo", () => {
  it("labels WMO codes in Croatian with a day or night icon", () => {
    expect(weatherInfo(0, true)).toEqual({ label: "Vedro", icon: "clear-day", precipitation: "none" });
    expect(weatherInfo(0, false).icon).toBe("clear-night");
    expect(weatherInfo(3, true)).toMatchObject({ label: "Oblačno", icon: "overcast-day" });
    expect(weatherInfo(63, true)).toMatchObject({ label: "Kiša", icon: "rain", precipitation: "rain" });
    expect(weatherInfo(73, false)).toMatchObject({ label: "Snijeg", icon: "snow", precipitation: "snow" });
    expect(weatherInfo(95, true)).toMatchObject({ label: "Grmljavina", precipitation: "storm" });
    expect(weatherInfo(51, true).precipitation).toBe("drizzle");
  });

  it("falls back for an unknown code", () => {
    expect(weatherInfo(42, true)).toEqual({ label: "", icon: "not-available", precipitation: "none" });
  });
});

describe("rainHint", () => {
  const start = new Date("2026-09-29T12:00:00+02:00");
  const hours = (wet: number[], code = 61): HourlyForecast[] =>
    Array.from({ length: 24 }, (_, i) => ({
      time: new Date(start.getTime() + i * 3_600_000).toISOString(),
      temperature: 15,
      apparentTemperature: 14,
      precipitationProbability: wet.includes(i) ? 80 : 5,
      precipitation: wet.includes(i) ? 1.2 : 0,
      weatherCode: wet.includes(i) ? code : 2,
      windSpeed: 10,
      windDirection: 180,
      isDay: true,
    }));
  const now = new Date("2026-09-29T12:20:00+02:00");

  it("says when rain starts", () => {
    expect(rainHint(hours([4, 5]), now, tz)).toBe("Kiša od 16 h");
  });

  it("says when ongoing rain stops", () => {
    expect(rainHint(hours([0, 1, 2]), now, tz)).toBe("Kiša do 15 h");
  });

  it("names snow and storms", () => {
    expect(rainHint(hours([3], 73), now, tz)).toBe("Snijeg od 15 h");
    expect(rainHint(hours([3], 95), now, tz)).toBe("Grmljavina od 15 h");
  });

  it("says when there is nothing for 12 hours", () => {
    expect(rainHint(hours([13]), now, tz)).toBe("Bez oborina idućih 12 h");
  });

  it("says rain continues when it doesn't stop within 12 hours", () => {
    expect(rainHint(hours([...Array(24).keys()]), now, tz)).toBe("Kiša i dalje");
  });

  it("reads hours on the location's clock", () => {
    expect(rainHint(hours([4, 5]), now, "UTC")).toBe("Kiša od 14 h");
  });
});
