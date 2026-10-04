export type AQILevel = "good" | "fair" | "moderate" | "poor" | "very_poor" | "extremely_poor";
export type PollenLevel = "none" | "low" | "moderate" | "high" | "very_high";
export type PollenType = "alder" | "birch" | "olive" | "grass" | "mugwort" | "ragweed";

/** Mirrors the backend's AirQualityResponse (app/models/airquality.go). */
export interface AirQuality {
  /** Current European Air Quality Index. */
  aqi: { value: number; level: AQILevel };
  pollen: {
    type: PollenType;
    /** grains/m³ */
    current: number;
    /** grains/m³ */
    todayMax: number;
    /** Of todayMax. */
    level: PollenLevel;
  }[];
  fetchedAt: string;
}
