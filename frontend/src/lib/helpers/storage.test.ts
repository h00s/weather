import { afterEach, describe, expect, it, vi } from "vitest";
import { load, save } from "./storage";

const isNumber = (v: unknown): v is number => typeof v === "number";

function fakeStorage(init: Record<string, string> = {}) {
  const data = new Map(Object.entries(init));
  return {
    getItem: (key: string) => data.get(key) ?? null,
    setItem: (key: string, value: string) => void data.set(key, value),
  };
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe("storage", () => {
  it("round-trips a valid value", () => {
    vi.stubGlobal("localStorage", fakeStorage());
    save("n", 42);
    expect(load("n", isNumber)).toBe(42);
  });

  it("treats a missing, corrupt or invalid value as nothing saved", () => {
    vi.stubGlobal("localStorage", fakeStorage({ corrupt: "{", wrong: '"text"' }));
    expect(load("missing", isNumber)).toBeUndefined();
    expect(load("corrupt", isNumber)).toBeUndefined();
    expect(load("wrong", isNumber)).toBeUndefined();
  });

  it("survives storage that throws or doesn't exist", () => {
    vi.stubGlobal("localStorage", {
      getItem: () => {
        throw new DOMException("blocked", "SecurityError");
      },
      setItem: () => {
        throw new DOMException("full", "QuotaExceededError");
      },
    });
    expect(load("n", isNumber)).toBeUndefined();
    expect(() => save("n", 1)).not.toThrow();

    vi.stubGlobal("localStorage", undefined);
    expect(load("n", isNumber)).toBeUndefined();
    expect(() => save("n", 1)).not.toThrow();
  });
});
