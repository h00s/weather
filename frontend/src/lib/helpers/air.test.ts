import { describe, expect, it } from "vitest";
import type { AirQuality, PollenLevel, PollenType } from "$lib/types/airquality";
import { activePollen, aqiLabels, pollenNames } from "./air";

const aq = (pollen: [PollenType, PollenLevel][]): AirQuality => ({
  aqi: { value: 30, level: "fair" },
  pollen: pollen.map(([type, level]) => ({ type, level, current: 0, todayMax: 0 })),
  fetchedAt: "",
});

describe("activePollen", () => {
  it("keeps the pollen in the air, the highest first", () => {
    const got = activePollen(aq([["grass", "low"], ["birch", "none"], ["ragweed", "very_high"], ["mugwort", "moderate"]]));
    expect(got.map((p) => p.type)).toEqual(["ragweed", "mugwort", "grass"]);
  });

  it("is empty outside the season", () => {
    expect(activePollen(aq([["grass", "none"], ["ragweed", "none"]]))).toEqual([]);
  });
});

describe("labels", () => {
  it("speak Croatian", () => {
    expect(aqiLabels.fair).toBe("Prihvatljiva");
    expect(pollenNames.ragweed).toBe("Ambrozija");
  });
});
