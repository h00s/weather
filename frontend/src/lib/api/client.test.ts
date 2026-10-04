import { describe, expect, it, vi } from "vitest";
import { api, ApiError, isAbortError } from "./client";

const respond = (status: number, body?: unknown) =>
  vi.fn<typeof globalThis.fetch>(
    async () =>
      new Response(body === undefined ? null : JSON.stringify(body), {
        status,
        headers: { "Content-Type": "application/json" },
      }),
  );

describe("api.get", () => {
  it("calls /api/v1 with the query, dropping empty values", async () => {
    const fetch = respond(200, { ok: true });

    const got = await api.get("/forecast", { fetch, query: { lat: 45.59, lon: 17.23, q: null, x: undefined } });

    expect(got).toEqual({ ok: true });
    expect(fetch.mock.calls[0][0]).toBe("/api/v1/forecast?lat=45.59&lon=17.23");
  });

  it("throws an ApiError carrying the status and Raptor's envelope", async () => {
    const fetch = respond(502, { code: 502, message: "Open-Meteo is unavailable" });

    const error = await api.get("/forecast", { fetch }).catch((e: unknown) => e);

    expect(error).toBeInstanceOf(ApiError);
    expect((error as ApiError).status).toBe(502);
    expect((error as ApiError).body).toEqual({ code: 502, message: "Open-Meteo is unavailable" });
  });

  it("recognizes the caller's own abort", () => {
    expect(isAbortError(new DOMException("aborted", "AbortError"))).toBe(true);
    expect(isAbortError(new Error("network"))).toBe(false);
  });
});
