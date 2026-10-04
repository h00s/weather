import { api } from "$lib/api/client";
import type { AirQuality } from "$lib/types/airquality";
import type { Coordinates } from "$lib/types/location";

export const fetchAirQuality = (at: Coordinates, signal?: AbortSignal, fetch?: typeof globalThis.fetch) =>
  api.get<AirQuality>("/air-quality", { query: { lat: at.latitude, lon: at.longitude }, signal, fetch });
