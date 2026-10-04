import { describe, expect, it } from "vitest";
import { fadeMask, scrollEdges, scrollStep } from "./scroll";

describe("scrollEdges", () => {
  it("says on which sides there is more to scroll to", () => {
    expect(scrollEdges(0, 400, 1344)).toEqual({ start: false, end: true });
    expect(scrollEdges(300, 400, 1344)).toEqual({ start: true, end: true });
    expect(scrollEdges(944, 400, 1344)).toEqual({ start: true, end: false });
  });

  it("ignores sub-pixel leftovers at either end", () => {
    expect(scrollEdges(0.6, 400, 1344).start).toBe(false);
    expect(scrollEdges(943.4, 400, 1344).end).toBe(false);
  });

  it("finds nothing to scroll when everything fits", () => {
    expect(scrollEdges(0, 400, 400)).toEqual({ start: false, end: false });
  });
});

describe("scrollStep", () => {
  it("moves most of a view, keeping a column for context", () => {
    expect(scrollStep(400)).toBe(320);
    expect(scrollStep(100)).toBe(80);
  });
});

describe("fadeMask", () => {
  it("fades only the sides that have more", () => {
    expect(fadeMask({ start: false, end: false })).toBe("none");
    expect(fadeMask({ start: false, end: true })).toBe("linear-gradient(to right, #000 0, #000 calc(100% - 3rem), transparent 100%)");
    expect(fadeMask({ start: true, end: true })).toBe("linear-gradient(to right, transparent 0, #000 3rem, #000 calc(100% - 3rem), transparent 100%)");
    expect(fadeMask({ start: true, end: false })).toBe("linear-gradient(to right, transparent 0, #000 3rem, #000 100%)");
  });
});
