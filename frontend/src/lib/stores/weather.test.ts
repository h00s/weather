import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import type { Forecast } from "$lib/types/forecast";

const services = vi.hoisted(() => ({ forecast: vi.fn(), airQuality: vi.fn(), warnings: vi.fn() }));
vi.mock("$lib/services/forecast", () => ({ fetchForecast: services.forecast }));
vi.mock("$lib/services/airquality", () => ({ fetchAirQuality: services.airQuality }));
vi.mock("$lib/services/warnings", () => ({ fetchWarnings: services.warnings }));

const daruvar = { latitude: 45.59, longitude: 17.23 };
const forecast = (fetchedAt: string) =>
  ({ current: { temperature: 12 }, hourly: [], daily: [], timezone: "Europe/Zagreb", elevation: 161, fetchedAt }) as unknown as Forecast;
const pending = () => new Promise<never>(() => {});
const settle = () => new Promise((resolve) => setTimeout(resolve, 0));

/** A fresh store per test: it is a module singleton. */
async function store() {
  vi.resetModules();
  return (await import("./weather.svelte")).weatherStore;
}

beforeEach(() => {
  for (const fn of Object.values(services)) fn.mockReset();
  services.airQuality.mockResolvedValue({ aqi: { value: 30, level: "fair" }, pollen: [], fetchedAt: "" });
  services.warnings.mockResolvedValue([]);
});

afterEach(() => {
  vi.useRealTimers();
  vi.unstubAllGlobals();
});

describe("weatherStore", () => {
  it("shows the forecast as soon as it arrives, whatever warnings and air quality are doing", async () => {
    services.forecast.mockResolvedValue(forecast(new Date().toISOString()));
    services.airQuality.mockImplementation(pending);
    services.warnings.mockImplementation(pending);
    const weather = await store();

    weather.show(daruvar);
    await settle();

    expect(weather.forecast?.current.temperature).toBe(12);
    expect(weather.loading).toBe(false);
  });

  it("dates the forecast by when the backend fetched it, and calls an old copy stale", async () => {
    const fetchedAt = new Date(Date.now() - 2 * 3_600_000).toISOString(); // served stale: Open-Meteo is down
    services.forecast.mockResolvedValue(forecast(fetchedAt));
    const weather = await store();

    weather.show(daruvar);
    await settle();

    expect(weather.updatedAt?.toISOString()).toBe(fetchedAt);
    expect(weather.stale).toBe(true);
  });

  it("calls a recent forecast fresh", async () => {
    services.forecast.mockResolvedValue(forecast(new Date(Date.now() - 5 * 60_000).toISOString()));
    const weather = await store();

    weather.show(daruvar);
    await settle();

    expect(weather.stale).toBe(false);
  });

  // Guard: an old fetchedAt must not make the poll fire again at once (it is timed from the last
  // request, not from the data's age).
  it("polls every 15 minutes even while the backend serves an old copy", async () => {
    vi.useFakeTimers();
    const listeners = { addEventListener: () => {}, removeEventListener: () => {} };
    vi.stubGlobal("document", { visibilityState: "visible", ...listeners });
    vi.stubGlobal("window", listeners);
    services.forecast.mockResolvedValue(forecast(new Date(Date.now() - 2 * 3_600_000).toISOString()));
    const weather = await store();

    const unsubscribe = weather.subscribe();
    weather.show(daruvar);
    await vi.advanceTimersByTimeAsync(60_000);
    expect(services.forecast).toHaveBeenCalledTimes(1);

    await vi.advanceTimersByTimeAsync(15 * 60_000);
    expect(services.forecast).toHaveBeenCalledTimes(2);
    unsubscribe();
  });
});
