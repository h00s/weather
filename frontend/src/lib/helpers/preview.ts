import { minutesOfDay } from "./format";

/** Design preview: `?at=19:40&code=61&cloud=90` shows the page at another time of day and in other
 *  weather, without waiting for it. */
export interface Preview {
  /** Minutes since midnight to show. */
  at?: number;
  /** A WMO weather code to show. */
  code?: number;
  /** Cloud cover in percent to show. */
  cloud?: number;
}

export function parsePreview(params: URLSearchParams): Preview {
  const preview: Preview = {};
  const at = /^(\d{1,2}):(\d{2})$/.exec(params.get("at") ?? "");
  if (at && Number(at[1]) < 24 && Number(at[2]) < 60) preview.at = Number(at[1]) * 60 + Number(at[2]);

  const code = Number.parseInt(params.get("code") ?? "", 10);
  if (Number.isInteger(code) && code >= 0 && code < 100) preview.code = code;

  const cloud = Number(params.get("cloud") ?? Number.NaN);
  if (params.has("cloud") && cloud >= 0 && cloud <= 100) preview.cloud = cloud;
  return preview;
}

/** Milliseconds to add to now so the location's clock reads `at` (minutes since midnight) today. */
export function previewOffset(at: number, now: Date, timeZone: string): number {
  return (at - minutesOfDay(now, timeZone)) * 60_000;
}
