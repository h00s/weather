import type { Warning } from "$lib/types/warning";
import { dayLabel, formatHour, isoDateOf } from "./format";

/** Rounded to the nearest hour, with the day when it isn't today: ["sutra", "6 h"]. */
function when(iso: string, now: Date, timeZone: string): [day: string, hour: string] {
  const t = new Date(Math.round(Date.parse(iso) / 3_600_000) * 3_600_000);
  const day = dayLabel(isoDateOf(t, timeZone), now, timeZone);
  return [day === "Danas" ? "" : day.toLowerCase(), formatHour(t, timeZone)];
}

/** When a warning applies: until when if it is in force, from when if it is still to come. */
export function warningWhen(w: Warning, now: Date, timeZone: string): string {
  if (Date.parse(w.onset) <= now.getTime()) {
    const [day, hour] = when(w.expires, now, timeZone);
    return day ? `do ${day} ${hour}` : `do ${hour}`;
  }
  const [day, hour] = when(w.onset, now, timeZone);
  return day ? `${day} od ${hour}` : `od ${hour}`;
}

/** A key for one warning in a list: every field the backend dedupes on, since two warnings may
 *  share an event and onset. */
export const warningKey = (w: Warning) => [w.level, w.type, w.event, w.description, w.onset, w.expires].join("|");
