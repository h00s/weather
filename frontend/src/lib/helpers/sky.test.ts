import { describe, expect, it } from "vitest";
import { contrastWithWhite, oklchCss, oklchToHex, skyAt } from "./sky";

const at = (hhmm: string) => new Date(`2026-09-29T${hhmm}:00+02:00`);
const sunrise = at("06:46");
const sunset = at("18:35");

describe("skyAt", () => {
  const clear = (hhmm: string) => skyAt(at(hhmm), sunrise, sunset, 0, "none");

  it("is darker at night than by day", () => {
    expect(clear("02:00").top.l).toBeLessThan(clear("12:00").top.l - 0.2);
  });

  it("stays blue by day and warms the horizon only around sunrise and sunset", () => {
    const hue = (h: number) => (h > 180 ? h - 360 : h); // warm hues sit near 0–90
    expect(hue(clear("06:30").bottom.h)).toBeLessThan(90);
    expect(hue(clear("18:40").bottom.h)).toBeLessThan(90);
    expect(clear("12:00").top.h).toBeGreaterThan(200);
    expect(clear("12:00").bottom.h).toBeGreaterThan(200);
    expect(clear("02:00").bottom.h).toBeGreaterThan(200);
  });

  it("changes gradually, not in steps", () => {
    const a = clear("06:40").bottom;
    const b = clear("06:41").bottom;
    expect(Math.abs(a.l - b.l)).toBeLessThan(0.02);
  });

  it("greys out under cloud and darkens in rain", () => {
    const sunny = skyAt(at("12:00"), sunrise, sunset, 0, "none");
    const overcast = skyAt(at("12:00"), sunrise, sunset, 100, "none");
    const rainy = skyAt(at("12:00"), sunrise, sunset, 100, "rain");
    expect(overcast.top.c).toBeLessThan(sunny.top.c / 2);
    expect(rainy.top.l).toBeLessThan(overcast.top.l);
  });

  it("formats as CSS", () => {
    expect(oklchCss({ l: 0.6, c: 0.13, h: 240 })).toBe("oklch(0.600 0.130 240.0)");
  });
});

describe("oklchToHex", () => {
  it("converts back to the brand navy", () => {
    const hex = oklchToHex({ l: 0.278, c: 0.073, h: 247.7 }); // #012a4a
    const [r, g, b] = [1, 3, 5].map((i) => Number.parseInt(hex.slice(i, i + 2), 16));
    expect(hex).toMatch(/^#[0-9a-f]{6}$/);
    expect(Math.abs(r - 0x01)).toBeLessThanOrEqual(2);
    expect(Math.abs(g - 0x2a)).toBeLessThanOrEqual(2);
    expect(Math.abs(b - 0x4a)).toBeLessThanOrEqual(2);
  });
});

describe("contrast with white text", () => {
  it("measures WCAG contrast", () => {
    expect(contrastWithWhite({ l: 0, c: 0, h: 0 })).toBeCloseTo(21, 0);
    expect(contrastWithWhite({ l: 1, c: 0, h: 0 })).toBeCloseTo(1, 1);
  });

  it("stays at 4.5:1 or more against every sky of the day, in any weather", () => {
    const worst = { contrast: Infinity, at: "" };
    for (let m = 0; m < 24 * 60; m += 5) {
      const now = new Date(at("00:00").getTime() + m * 60_000);
      for (const cloud of [0, 50, 100]) {
        for (const p of ["none", "snow"] as const) {
          const sky = skyAt(now, sunrise, sunset, cloud, p);
          for (const c of [sky.top, sky.bottom]) {
            const contrast = contrastWithWhite(c);
            if (contrast < worst.contrast) Object.assign(worst, { contrast, at: `${m} min, cloud ${cloud}, ${p}` });
          }
        }
      }
    }
    expect(worst.contrast, worst.at).toBeGreaterThanOrEqual(4.5);
  });

  it("keeps the daytime sky colourful rather than muddy", () => {
    expect(skyAt(at("12:00"), sunrise, sunset, 0, "none").top.c).toBeGreaterThan(0.08);
  });
});
