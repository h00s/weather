import { describe, expect, it } from "vitest";
import { areaPath, extent, linePath, scale } from "./sparkline";

describe("extent", () => {
  it("ignores nulls", () => {
    expect(extent([3, null, -1, 7])).toEqual([-1, 7]);
    expect(extent([null, null])).toBeNull();
  });
});

describe("scale", () => {
  it("maps the domain onto the range", () => {
    const y = scale([0, 10], [50, 0]);
    expect(y(0)).toBe(50);
    expect(y(10)).toBe(0);
    expect(y(5)).toBe(25);
  });

  it("centres a flat domain", () => {
    expect(scale([4, 4], [50, 0])(4)).toBe(25);
  });
});

describe("linePath", () => {
  it("draws a line across the width", () => {
    expect(linePath([0, 5, 10], 100, 50, [0, 10])).toBe("M0,50L50,25L100,0");
  });

  it("breaks the line at gaps", () => {
    expect(linePath([0, null, 10, 10], 90, 50, [0, 10])).toBe("M0,50M60,0L90,0");
  });
});

describe("areaPath", () => {
  it("closes each run down to the baseline", () => {
    expect(areaPath([0, 10], 100, 50, [0, 10])).toBe("M0,50L100,0L100,50L0,50Z");
    expect(areaPath([0, null, 10, 10], 90, 50, [0, 10])).toBe("M60,0L90,0L90,50L60,50Z");
  });
});
