// Device storage is a convenience. Private browsing, blocked site data or a full quota make it throw
// or vanish, and an older app version may have left another shape behind, so every read is
// validated and every failure simply means "nothing saved".

/** The value saved under key; undefined when there is none, it doesn't parse, or it fails valid. */
export function load<T>(key: string, valid: (value: unknown) => value is T): T | undefined {
  try {
    const raw = globalThis.localStorage?.getItem(key);
    if (raw == null) return undefined;
    const value: unknown = JSON.parse(raw);
    return valid(value) ? value : undefined;
  } catch {
    return undefined;
  }
}

/** Saves value under key; does nothing when storage is unavailable or full. */
export function save(key: string, value: unknown): void {
  try {
    globalThis.localStorage?.setItem(key, JSON.stringify(value));
  } catch {
    // Best effort: the app works without it, it just won't remember.
  }
}
