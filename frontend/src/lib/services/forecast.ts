import { api } from "$lib/api/client";
import type { Forecast } from "$lib/types/forecast";
import type { Coordinates } from "$lib/types/location";

export const fetchForecast = (at: Coordinates, signal?: AbortSignal, fetch?: typeof globalThis.fetch) =>
  api.get<Forecast>("/forecast", { query: { lat: at.latitude, lon: at.longitude }, signal, fetch });
