import { describe, expect, it } from "vitest";
import { uvBand } from "./uv";

describe("uvBand", () => {
  it("follows the WHO bands on the rounded index", () => {
    expect([0, 2.4, 2.55, 5, 6, 8, 11].map((i) => uvBand(i).label)).toEqual([
      "Nizak",
      "Nizak",
      "Umjeren",
      "Umjeren",
      "Visok",
      "Vrlo visok",
      "Ekstreman",
    ]);
    expect(uvBand(9).level).toBe(3);
  });
});
