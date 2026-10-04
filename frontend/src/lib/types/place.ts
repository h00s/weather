/** Mirrors the backend's PlaceResponse (app/models/place.go). */
export interface Place {
  name: string;
  /** e.g. "Bjelovarsko-bilogorska županija"; "" when unknown */
  county: string;
  latitude: number;
  longitude: number;
}
