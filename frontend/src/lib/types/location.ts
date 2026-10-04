export interface Coordinates {
  latitude: number;
  longitude: number;
}

/** What the app shows weather for, as stored on this device. Coordinates are rounded to 2 decimals
 *  (~1 km) before they are stored or sent. */
export interface SavedLocation extends Coordinates {
  kind: "gps" | "place";
  /** A settlement; null for a position with no Croatian place nearby ("Moja lokacija"). */
  name: string | null;
  county: string | null;
}
