import { api } from "$lib/api/client";
import type { Coordinates } from "$lib/types/location";
import type { Place } from "$lib/types/place";

/** Croatian settlements by name (2 to 60 characters), best match first. */
export const searchPlaces = (q: string, signal?: AbortSignal, fetch?: typeof globalThis.fetch) =>
  api.get<Place[]>("/places", { query: { q }, signal, fetch });

/** The settlement at a point; rejects with an ApiError 404 abroad or at sea. */
export const nearestPlace = (at: Coordinates, signal?: AbortSignal, fetch?: typeof globalThis.fetch) =>
  api.get<Place>("/places/nearest", { query: { lat: at.latitude, lon: at.longitude }, signal, fetch });
