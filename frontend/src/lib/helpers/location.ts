import type { SavedLocation } from "$lib/types/location";
import { cellKey } from "./geo";

export const MAX_RECENT = 5;

/** A SavedLocation as stored by this version of the app. */
export function isSavedLocation(v: unknown): v is SavedLocation {
  if (typeof v !== "object" || v === null) return false;
  const l = v as Record<string, unknown>;
  return (
    (l.kind === "gps" || l.kind === "place") &&
    (l.name === null || typeof l.name === "string") &&
    (l.county === null || typeof l.county === "string") &&
    typeof l.latitude === "number" &&
    Math.abs(l.latitude) <= 90 &&
    typeof l.longitude === "number" &&
    Math.abs(l.longitude) <= 180
  );
}

export const isRecentList = (v: unknown): v is SavedLocation[] => Array.isArray(v) && v.every(isSavedLocation);

/** The recents with place in front: newest first, one entry per ~1 km cell, at most MAX_RECENT.
 *  GPS fixes aren't recents, since "Koristi moju lokaciju" is always offered on its own. */
export function addRecent(recent: SavedLocation[], place: SavedLocation): SavedLocation[] {
  if (place.kind !== "place") return recent;
  const key = cellKey(place);
  return [place, ...recent.filter((r) => cellKey(r) !== key)].slice(0, MAX_RECENT);
}

/** What the header calls a location. */
export const locationLabel = (l: SavedLocation) => l.name ?? "Moja lokacija";

/** Croatian text for a failed position request: GeolocationPositionError's code, or 0 when the
 *  browser can't locate at all. */
export function positionErrorMessage(code: number): string {
  switch (code) {
    case 1:
      return "Pristup lokaciji nije dopušten. Dopustite ga u postavkama preglednika ili odaberite mjesto.";
    case 2:
      return "Lokacija trenutno nije dostupna. Pokušajte ponovno ili odaberite mjesto.";
    case 3:
      return "Određivanje lokacije traje predugo. Pokušajte ponovno ili odaberite mjesto.";
    default:
      return "Ovaj preglednik ne može odrediti lokaciju. Odaberite mjesto.";
  }
}
