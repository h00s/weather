import { describe, expect, it } from "vitest";
import { parsePreview, previewOffset } from "./preview";

describe("parsePreview", () => {
  it("reads a time, a WMO code and a cloud cover", () => {
    expect(parsePreview(new URLSearchParams("at=19:40&code=61&cloud=90"))).toEqual({ at: 1180, code: 61, cloud: 90 });
  });

  it("ignores what doesn't parse", () => {
    expect(parsePreview(new URLSearchParams("at=25:00&code=abc&cloud=150"))).toEqual({});
    expect(parsePreview(new URLSearchParams(""))).toEqual({});
  });
});

describe("previewOffset", () => {
  it("moves now to the asked time on the location's clock", () => {
    expect(previewOffset(1180, new Date("2026-09-29T19:00:00+02:00"), "Europe/Zagreb")).toBe(40 * 60_000);
  });
});
