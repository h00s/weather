import { describe, expect, it } from "vitest";
import {
  dayLabel,
  formatAgo,
  formatClock,
  formatDegrees,
  formatDistance,
  formatDuration,
  formatHour,
  formatMillimetres,
  formatPercent,
  formatPressure,
  formatShortDate,
  formatSpeed,
  isoDateOf,
  minutesOfDay,
} from "./format";

const tz = "Europe/Zagreb";
const d = new Date("2026-09-29T07:05:00+02:00");

describe("times on the location's clock", () => {
  it("formats clock times and hours in the given zone", () => {
    expect(formatClock(d, tz)).toBe("07:05");
    expect(formatHour(d, tz)).toBe("7 h");
    expect(formatClock(d, "UTC")).toBe("05:05");
    expect(minutesOfDay(d, tz)).toBe(425);
  });

  it("dates by the location's calendar, not the device's", () => {
    expect(isoDateOf(new Date("2026-09-29T23:30:00Z"), tz)).toBe("2026-09-30");
  });

  it("names days relative to today, then by weekday and date", () => {
    expect(dayLabel("2026-09-29", d, tz)).toBe("Danas");
    expect(dayLabel("2026-09-30", d, tz)).toBe("Sutra");
    expect(dayLabel("2026-10-01", d, tz)).toBe("Čet 1. 10.");
    expect(dayLabel("2026-10-04", d, tz)).toBe("Ned 4. 10.");
  });

  it("writes a short date without zero padding", () => {
    expect(formatShortDate("2026-10-04")).toBe("4. 10.");
    expect(formatShortDate("2026-12-31")).toBe("31. 12.");
  });
});

describe("values", () => {
  it("rounds degrees, never showing -0", () => {
    expect(formatDegrees(17.46)).toBe("17°");
    expect(formatDegrees(-0.4)).toBe("0°");
    expect(formatDegrees(null)).toBe("–");
  });

  it("formats measurements the Croatian way", () => {
    expect(formatPercent(47.6)).toBe("48 %");
    expect(formatMillimetres(4.25)).toBe("4,3 mm");
    expect(formatMillimetres(0)).toBe("0 mm");
    expect(formatSpeed(12.4)).toBe("12 km/h");
    expect(formatPressure(1013.2)).toBe("1.013 hPa");
    expect(formatDistance(41_300)).toBe("41 km");
    expect(formatDistance(2_500)).toBe("2,5 km");
    expect(formatDistance(820)).toBe("800 m");
    expect(formatDuration(41_590)).toBe("11 h 33 min");
  });

  it("says how long ago", () => {
    const now = new Date("2026-09-29T12:00:00Z");
    const ago = (ms: number) => formatAgo(new Date(now.getTime() - ms), now);
    expect(ago(20_000)).toBe("upravo");
    expect(ago(5 * 60_000)).toBe("prije 5 min");
    expect(ago(125 * 60_000)).toBe("prije 2 h");
    expect(ago(24 * 3_600_000)).toBe("prije 1 dan");
    expect(ago(3 * 24 * 3_600_000)).toBe("prije 3 dana");
  });
});
