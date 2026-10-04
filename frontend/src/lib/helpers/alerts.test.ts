import { describe, expect, it } from "vitest";
import type { Warning } from "$lib/types/warning";
import { warningKey, warningWhen } from "./alerts";

describe("warningWhen", () => {
  const tz = "Europe/Zagreb";
  const w = (onset: string, expires: string) => ({ onset, expires }) as Warning;
  const now = new Date("2026-09-29T12:00:00+02:00");

  it("says until when a warning in force lasts", () => {
    expect(warningWhen(w("2026-09-29T10:00:00+02:00", "2026-09-29T17:59:59+02:00"), now, tz)).toBe("do 18 h");
  });

  it("says from when an upcoming warning starts, naming tomorrow", () => {
    expect(warningWhen(w("2026-09-29T16:00:00+02:00", "2026-09-29T22:00:00+02:00"), now, tz)).toBe("od 16 h");
    expect(warningWhen(w("2026-09-30T06:00:00+02:00", "2026-09-30T12:00:00+02:00"), now, tz)).toBe("sutra od 6 h");
  });

  it("names the day when a warning in force ends tomorrow", () => {
    expect(warningWhen(w("2026-09-29T10:00:00+02:00", "2026-09-30T05:59:59+02:00"), now, tz)).toBe("do sutra 6 h");
  });

  it("names a later day by weekday and date", () => {
    expect(warningWhen(w("2026-10-01T06:00:00+02:00", "2026-10-01T12:00:00+02:00"), now, tz)).toBe("čet 1. 10. od 6 h");
  });
});

describe("warningKey", () => {
  // The backend drops only exact duplicates, so two warnings can share an event and onset; Svelte
  // throws on a repeated each-key, even in production, and the page would not render.
  it("tells apart warnings that differ in any field the backend keeps", () => {
    const base: Warning = {
      level: "yellow",
      type: "wind",
      event: "Žuto upozorenje za vjetar",
      description: "Jaka bura.",
      onset: "2026-09-25T09:00:00+02:00",
      expires: "2026-09-25T15:00:00+02:00",
    };
    const keys = [
      base,
      { ...base, description: "Jaka bura s olujnim udarima." },
      { ...base, expires: "2026-09-25T18:00:00+02:00" },
      { ...base, level: "orange" as const },
      { ...base, type: "rain" },
    ].map(warningKey);
    expect(new Set(keys).size).toBe(keys.length);
    expect(warningKey({ ...base })).toBe(warningKey(base));
  });
});
