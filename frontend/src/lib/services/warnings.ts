import { api } from "$lib/api/client";
import type { Coordinates } from "$lib/types/location";
import type { Warning } from "$lib/types/warning";

export const fetchWarnings = (at: Coordinates, signal?: AbortSignal, fetch?: typeof globalThis.fetch) =>
  api.get<Warning[]>("/warnings", { query: { lat: at.latitude, lon: at.longitude }, signal, fetch });
