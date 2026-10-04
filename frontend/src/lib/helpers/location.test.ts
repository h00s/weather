import { describe, expect, it } from "vitest";
import { cities } from "$lib/data/cities";
import type { SavedLocation } from "$lib/types/location";
import { addRecent, isSavedLocation, locationLabel, MAX_RECENT, positionErrorMessage } from "./location";

const place = (name: string, latitude: number, longitude: number): SavedLocation => ({
  kind: "place",
  name,
  county: null,
  latitude,
  longitude,
});

describe("isSavedLocation", () => {
  it("accepts what the app stores", () => {
    expect(isSavedLocation(place("Daruvar", 45.59, 17.22))).toBe(true);
    expect(isSavedLocation({ kind: "gps", name: null, county: null, latitude: 43.86, longitude: 18.41 })).toBe(true);
  });

  it("rejects anything else, so a corrupt or older value means nothing saved", () => {
    const bad = [
      null,
      "Daruvar",
      42,
      {},
      { ...place("x", 1, 2), kind: "city" },
      { ...place("x", 1, 2), latitude: "45" },
      place("x", 91, 0),
      place("x", 0, 181),
      { kind: "place", name: "x", latitude: 1, longitude: 2 },
    ];
    for (const v of bad) expect(isSavedLocation(v), JSON.stringify(v)).toBe(false);
  });
});

describe("addRecent", () => {
  it("puts the newest first, once per ~1 km cell, at most MAX_RECENT", () => {
    let recent: SavedLocation[] = [];
    for (let i = 0; i < 7; i++) recent = addRecent(recent, place(`P${i}`, 45 + i / 10, 16));
    expect(recent.map((r) => r.name)).toEqual(["P6", "P5", "P4", "P3", "P2"]);
    expect(recent).toHaveLength(MAX_RECENT);

    recent = addRecent(recent, place("P4 again", 45.4012, 16.0004));
    expect(recent.map((r) => r.name)).toEqual(["P4 again", "P6", "P5", "P3", "P2"]);
  });

  it("leaves GPS fixes out", () => {
    const recent = [place("Daruvar", 45.59, 17.22)];
    expect(addRecent(recent, { kind: "gps", name: "Split", county: null, latitude: 43.51, longitude: 16.44 })).toBe(recent);
  });
});

describe("locationLabel", () => {
  it("calls a position with no Croatian place nearby 'Moja lokacija'", () => {
    expect(locationLabel({ kind: "gps", name: null, county: null, latitude: 43.86, longitude: 18.41 })).toBe("Moja lokacija");
    expect(locationLabel(place("Daruvar", 45.59, 17.22))).toBe("Daruvar");
  });
});

describe("positionErrorMessage", () => {
  it("explains each failure in Croatian", () => {
    const messages = [0, 1, 2, 3].map(positionErrorMessage);
    expect(new Set(messages).size).toBe(4);
    expect(messages[1]).toContain("nije dopušten");
    expect(messages[3]).toContain("predugo");
  });
});

describe("cities", () => {
  it("offers 20 distinct, valid, rounded places", () => {
    expect(cities).toHaveLength(20);
    expect(new Set(cities.map((c) => c.name)).size).toBe(20);
    for (const c of cities) {
      expect(isSavedLocation(c), c.name ?? "").toBe(true);
      expect(Math.round(c.latitude * 100) / 100).toBe(c.latitude);
      expect(Math.round(c.longitude * 100) / 100).toBe(c.longitude);
    }
  });
});
