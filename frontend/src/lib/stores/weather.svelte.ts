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
/** A forecast the backend fetched longer ago than this is shown as stale: it is serving its last
 *  good copy while Open-Meteo is down. */
const OLD_DATA_MS = 30 * 60_000;

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
  /** When the last refresh started. Polls are timed from it, not from the data's age, which an
   *  outage keeps old. */
  let lastAttempt = 0;

  /** The next poll, POLL_MS after the last refresh; none while hidden or unwatched. */
  function schedule() {
    if (timer !== null) clearTimeout(timer);
    timer = null;
    if (consumers === 0 || document.visibilityState === "hidden") return;
    timer = setTimeout(() => void refresh(), Math.max(0, lastAttempt + POLL_MS - Date.now()));
  }

  /** Fetches forecast, air quality and warnings for the current cell in parallel. Each lands as
   *  soon as it arrives, so a slow or failing Meteoalarm never holds back or blanks the forecast. A
   *  newer refresh, a location change or the last consumer leaving aborts this one, and its late
   *  results are dropped. */
  async function refresh(): Promise<void> {
    const target = untrack(() => at);
    if (!target) return;
    controller?.abort();
    const mine = (controller = new AbortController());
    const current = () => controller === mine;
    lastAttempt = Date.now();
    loading = true;

    const forecastLanded = fetchForecast(target, mine.signal)
      .then(
        (f) => {
          if (!current()) return;
          forecast = f;
          updatedAt = new Date(f.fetchedAt);
          error = null;
          save(CACHE_KEY, {
            version: CACHE_VERSION,
            cell: cellKey(target),
            savedAt: new Date().toISOString(),
            forecast: f,
          } satisfies CachedForecast);
        },
        (e: unknown) => {
          if (current() && !isAbortError(e)) error = e;
        },
      )
      .finally(() => {
        if (current()) loading = false;
      });
    // A failure here leaves the card or banner out; the forecast stands on its own.
    const othersLanded = [
      fetchAirQuality(target, mine.signal).then(
        (a) => {
          if (current()) airQuality = a;
        },
        () => {},
      ),
      fetchWarnings(target, mine.signal).then(
        (w) => {
          if (current()) warnings = w;
        },
        () => {},
      ),
    ];

    await Promise.allSettled([forecastLanded, ...othersLanded]);
    if (!current()) return;
    controller = null;
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
      updatedAt = cached ? new Date(cached.forecast.fetchedAt) : null;
      airQuality = undefined;
      warnings = undefined;
      error = null;
      void refresh();
    });
  }

  function onVisibility() {
    if (document.visibilityState === "visible" && Date.now() - lastAttempt > STALE_MS) void refresh();
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
    /** When the backend fetched the forecast shown from Open-Meteo. */
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
    /** There is a forecast, but it may be out of date: the last refresh failed, or the backend
     *  could only serve an old copy. Not while a refresh is under way. */
    get stale() {
      if (forecast === undefined || loading) return false;
      return error !== null || (updatedAt !== null && Date.now() - updatedAt.getTime() > OLD_DATA_MS);
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
