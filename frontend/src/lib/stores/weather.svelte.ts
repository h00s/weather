import { untrack } from "svelte";
import { isAbortError } from "$lib/api/client";
import { CACHE_VERSION, isCachedForecast, usableCache, type CachedForecast } from "$lib/helpers/forecast";
import { cellKey } from "$lib/helpers/geo";
import { load, save } from "$lib/helpers/storage";
import { fetchAirQuality } from "$lib/services/airquality";
import { fetchForecast } from "$lib/services/forecast";
import { fetchWarnings } from "$lib/services/warnings";
import type { AirQuality } from "$lib/types/airquality";
import type { Forecast } from "$lib/types/forecast";
import type { Coordinates } from "$lib/types/location";
import type { Warning } from "$lib/types/warning";

const CACHE_KEY = "vrijeme.forecast";
/** The backend caches the forecast for 10 minutes; asking every 15 keeps an open page current. */
const POLL_MS = 15 * 60_000;
/** Coming back to the page after this long refreshes at once instead of at the next poll. */
const STALE_MS = 10 * 60_000;

function createWeatherStore() {
  let at = $state<Coordinates | null>(null);
  let forecast = $state<Forecast>();
  let airQuality = $state<AirQuality>();
  let warnings = $state<Warning[]>();
  let updatedAt = $state<Date | null>(null);
  let error = $state<unknown>(null);
  let loading = $state(false);

  let controller: AbortController | null = null;
  let timer: ReturnType<typeof setTimeout> | null = null;
  let consumers = 0;

  const age = () => {
    const since = untrack(() => updatedAt);
    return since ? Date.now() - since.getTime() : Infinity;
  };

  /** The next poll, POLL_MS after the last update; none while hidden or unwatched. */
  function schedule() {
    if (timer !== null) clearTimeout(timer);
    timer = null;
    if (consumers === 0 || document.visibilityState === "hidden") return;
    timer = setTimeout(() => void refresh(), Math.max(0, POLL_MS - age()));
  }

  /** Fetches forecast, air quality and warnings for the current cell in parallel. Each lands on its
   *  own, so a Meteoalarm outage never blanks the forecast. A newer refresh, a location change or
   *  the last consumer leaving aborts this one, and its results are dropped. */
  async function refresh(): Promise<void> {
    const target = untrack(() => at);
    if (!target) return;
    controller?.abort();
    const mine = (controller = new AbortController());
    loading = true;

    const [f, a, w] = await Promise.allSettled([
      fetchForecast(target, mine.signal),
      fetchAirQuality(target, mine.signal),
      fetchWarnings(target, mine.signal),
    ]);
    if (controller !== mine) return; // superseded: a newer refresh owns the state now
    controller = null;
    loading = false;

    if (f.status === "fulfilled") {
      forecast = f.value;
      updatedAt = new Date();
      error = null;
      save(CACHE_KEY, {
        version: CACHE_VERSION,
        cell: cellKey(target),
        savedAt: updatedAt.toISOString(),
        forecast: f.value,
      } satisfies CachedForecast);
    } else if (!isAbortError(f.reason)) {
      error = f.reason;
    }
    if (a.status === "fulfilled") airQuality = a.value;
    if (w.status === "fulfilled") warnings = w.value;
    schedule();
  }

  /** Switches to location's ~1 km cell: paints the forecast saved on the device for it, when recent,
   *  then refreshes. The same cell again changes nothing. */
  function show(location: Coordinates) {
    untrack(() => {
      const cell = cellKey(location);
      if (at && cellKey(at) === cell) return;
      at = { latitude: location.latitude, longitude: location.longitude };
      const cached = usableCache(load(CACHE_KEY, isCachedForecast), cell, new Date());
      forecast = cached?.forecast;
      updatedAt = cached ? new Date(cached.savedAt) : null;
      airQuality = undefined;
      warnings = undefined;
      error = null;
      void refresh();
    });
  }

  function onVisibility() {
    if (document.visibilityState === "visible" && age() > STALE_MS) void refresh();
    else schedule();
  }

  function onOnline() {
    void refresh();
  }

  return {
    get forecast() {
      return forecast;
    },
    get airQuality() {
      return airQuality;
    },
    get warnings() {
      return warnings;
    },
    /** When the forecast shown was fetched; a cached one keeps its own time. */
    get updatedAt() {
      return updatedAt;
    },
    /** The last forecast refresh's failure, cleared by the next success. */
    get error() {
      return error;
    },
    get loading() {
      return loading;
    },
    /** There is a forecast, but the last refresh failed: what's shown may be out of date. */
    get stale() {
      return forecast !== undefined && error !== null;
    },
    show,
    refresh,
    /** `$effect(() => weatherStore.subscribe())` polls while the page is mounted and visible. */
    subscribe(): () => void {
      return untrack(() => {
        if (consumers++ === 0) {
          document.addEventListener("visibilitychange", onVisibility);
          window.addEventListener("online", onOnline);
          schedule();
        }
        return () => {
          if (--consumers > 0) return;
          document.removeEventListener("visibilitychange", onVisibility);
          window.removeEventListener("online", onOnline);
          if (timer !== null) clearTimeout(timer);
          timer = null;
          controller?.abort();
          controller = null;
          loading = false;
        };
      });
    },
  };
}

export const weatherStore = createWeatherStore();
