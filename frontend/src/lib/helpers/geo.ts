import type { Coordinates } from "$lib/types/location";

const round2 = (v: number) => Math.round(v * 100) / 100 || 0; // || 0: never -0

/** Rounded to 2 decimals (~1 km), the backend's cache cell: the only precision stored on the
 *  device or sent to the API. */
export const roundCoordinates = ({ latitude, longitude }: Coordinates): Coordinates => ({
  latitude: round2(latitude),
  longitude: round2(longitude),
});

/** The ~1 km cell, e.g. "45.59,17.23": two locations in one cell get the same forecast. */
export function cellKey(c: Coordinates): string {
  const r = roundCoordinates(c);
  return `${r.latitude.toFixed(2)},${r.longitude.toFixed(2)}`;
}
