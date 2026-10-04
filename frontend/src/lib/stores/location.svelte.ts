import { untrack } from "svelte";
import { ApiError } from "$lib/api/client";
import { cellKey, roundCoordinates } from "$lib/helpers/geo";
import { addRecent, isRecentList, isSavedLocation, positionErrorMessage } from "$lib/helpers/location";
import { load, save } from "$lib/helpers/storage";
import { nearestPlace } from "$lib/services/places";
import type { SavedLocation } from "$lib/types/location";

const CURRENT_KEY = "vrijeme.location";
const RECENT_KEY = "vrijeme.recent";

/** How old a position the device already has may be before it is asked again. */
const POSITION_MAX_AGE_MS = 15 * 60_000;

function currentPosition(): Promise<GeolocationPosition> {
  return new Promise((resolve, reject) =>
    navigator.geolocation.getCurrentPosition(resolve, reject, { timeout: 10_000, maximumAge: POSITION_MAX_AGE_MS }),
  );
}

function createLocationStore() {
  let current = $state<SavedLocation | null>(load(CURRENT_KEY, isSavedLocation) ?? null);
  let recent = $state<SavedLocation[]>(load(RECENT_KEY, isRecentList) ?? []);
  let locating = $state(false);
  let error = $state<string | null>(null);

  function set(location: SavedLocation) {
    current = location;
    save(CURRENT_KEY, location);
  }

  /** Asks the device where it is, names the spot after the nearest Croatian settlement, and shows
   *  weather there. Resolves false, with error set, when there is no position to be had. */
  async function locate(): Promise<boolean> {
    if (!("geolocation" in navigator)) {
      error = positionErrorMessage(0);
      return false;
    }
    locating = true;
    error = null;
    try {
      const { coords } = await currentPosition();
      const at = roundCoordinates({ latitude: coords.latitude, longitude: coords.longitude });
      const previous = untrack(() => current);
      let name: string | null = null;
      let county: string | null = null;
      try {
        const place = await nearestPlace(at);
        name = place.name;
        county = place.county || null;
      } catch (e) {
        // A 404 means abroad or at sea: "Moja lokacija". Anything else (offline) keeps the name
        // this cell already had.
        const abroad = e instanceof ApiError && e.status === 404;
        if (!abroad && previous && cellKey(previous) === cellKey(at)) {
          name = previous.name;
          county = previous.county;
        }
      }
      set({ kind: "gps", name, county, ...at });
      return true;
    } catch (e) {
      error = positionErrorMessage(e instanceof GeolocationPositionError ? e.code : 0);
      return false;
    } finally {
      locating = false;
    }
  }

  return {
    /** What the app shows weather for; null until the reader picks something (the welcome screen). */
    get current() {
      return current;
    },
    /** Places picked before, newest first. */
    get recent() {
      return recent;
    },
    /** A position request is running. */
    get locating() {
      return locating;
    },
    /** Why the last position request failed, in Croatian; cleared by the next attempt or pick. */
    get error() {
      return error;
    },

    /** Shows weather for a picked place and remembers it among the recents. */
    choose(place: SavedLocation) {
      const location = { ...place, ...roundCoordinates(place) };
      error = null;
      set(location);
      recent = addRecent(untrack(() => recent), location);
      save(RECENT_KEY, recent);
    },

    locate,

    /** In GPS mode, follows the device on open without a prompt, but only when permission is
     *  already granted: a returning reader never meets a cold permission dialog. */
    async relocateIfGranted(): Promise<void> {
      if (untrack(() => current)?.kind !== "gps") return;
      try {
        const status = await navigator.permissions?.query({ name: "geolocation" });
        if (status?.state === "granted") await locate();
      } catch {
        // No Permissions API: stay on the saved position until the reader refreshes.
      }
    },
  };
}

export const locationStore = createLocationStore();
