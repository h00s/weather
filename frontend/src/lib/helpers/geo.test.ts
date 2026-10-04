import { describe, expect, it } from "vitest";
import { cellKey, roundCoordinates } from "./geo";

describe("roundCoordinates", () => {
  it("rounds to the ~1 km cell, never to -0", () => {
    expect(roundCoordinates({ latitude: 45.5936, longitude: 17.2251 })).toEqual({ latitude: 45.59, longitude: 17.23 });
    const zero = roundCoordinates({ latitude: -0.001, longitude: 0.004 });
    expect(Object.is(zero.latitude, 0)).toBe(true);
  });
});

describe("cellKey", () => {
  it("is shared by neighbours in one cell", () => {
    expect(cellKey({ latitude: 45.5936, longitude: 17.2251 })).toBe("45.59,17.23");
    expect(cellKey({ latitude: 45.5912, longitude: 17.2263 })).toBe("45.59,17.23");
    expect(cellKey({ latitude: 45.6, longitude: 17 })).toBe("45.60,17.00");
  });
});
