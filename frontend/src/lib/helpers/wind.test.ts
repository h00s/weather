import { describe, expect, it } from "vitest";
import { compassName, compassPoint, windArrowDegrees, windStrength } from "./wind";

describe("compass", () => {
  it("names where the wind blows from, on 8 Croatian points", () => {
    expect([0, 22, 23, 90, 225, 359, -45, 405].map(compassPoint)).toEqual(["S", "S", "SI", "I", "JZ", "S", "SZ", "SI"]);
    expect(compassName(135)).toBe("jugoistok");
  });

  it("points the arrow where the wind blows to", () => {
    expect(windArrowDegrees(225)).toBe(45);
    expect(windArrowDegrees(0)).toBe(180);
  });
});

describe("windStrength", () => {
  it("groups Beaufort forces in Croatian", () => {
    expect([1, 12, 30, 50, 70, 120].map(windStrength)).toEqual([
      "Tiho",
      "Slab vjetar",
      "Umjeren vjetar",
      "Jak vjetar",
      "Olujni vjetar",
      "Orkanski vjetar",
    ]);
  });
});
