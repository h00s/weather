/** Raptor's error envelope: `{"code":502,"message":"Open-Meteo is unavailable"}`. `message` is
 *  English and meant for developers. */
export interface ApiErrorBody {
  code?: number;
  message?: string;
  attrs?: Record<string, unknown>;
}

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public body?: ApiErrorBody,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

export interface ApiOptions {
  /** SvelteKit's fetch, passed in from a load. */
  fetch?: typeof globalThis.fetch;
  /** Query params: null and undefined entries are dropped, the rest are encoded. */
  query?: Record<string, string | number | boolean | null | undefined>;
  /** Aborts the request. The rejection is an AbortError (see isAbortError). */
  signal?: AbortSignal;
}

// The app only reads, and has no session: this is the conventions' client without the 401 seam,
// writes, uploads and downloads.

/** Same origin everywhere: the Go server serves the SPA in production, and Vite proxies /api in
 *  development. So there is nothing to configure, no CORS, and no $env. */
const BASE_URL = "/api/v1";

function buildUrl(endpoint: string, query?: ApiOptions["query"]): string {
  const url = `${BASE_URL}${endpoint}`;
  if (!query) return url;
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(query)) {
    if (value !== undefined && value !== null) params.set(key, String(value));
  }
  const qs = params.toString();
  return qs ? `${url}?${qs}` : url;
}

async function request<T>(method: string, endpoint: string, opts: ApiOptions = {}): Promise<T> {
  const { fetch: customFetch, query, signal } = opts;

  const response = await (customFetch ?? fetch)(buildUrl(endpoint, query), { method, signal });

  if (!response.ok) {
    const errorBody: ApiErrorBody = await response.json().catch(() => ({}));
    throw new ApiError(errorBody.message ?? `Request failed (${response.status})`, response.status, errorBody);
  }
  if (response.status === 204 || response.headers.get("content-length") === "0") {
    return undefined as T;
  }
  return response.json() as Promise<T>;
}

/** True when a rejection is the caller's own abort rather than a failure. */
export const isAbortError = (e: unknown) => e instanceof DOMException && e.name === "AbortError";

export const api = {
  get: <T>(endpoint: string, opts?: ApiOptions) => request<T>("GET", endpoint, opts),
};
