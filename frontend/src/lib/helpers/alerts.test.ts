import { describe, expect, it } from "vitest";
import type { Warning } from "$lib/types/warning";
import { warningWhen } from "./alerts";

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
