// Times and dates are shown on the forecast location's clock and calendar: pass its IANA zone
// (Forecast.timezone), never the device's.

const dateFormats = new Map<string, Intl.DateTimeFormat>();

/** An Intl.DateTimeFormat, built once per locale and options: constructing one is slow. */
function dateFormat(locale: string, options: Intl.DateTimeFormatOptions): Intl.DateTimeFormat {
  const key = `${locale}|${JSON.stringify(options)}`;
  let format = dateFormats.get(key);
  if (!format) dateFormats.set(key, (format = new Intl.DateTimeFormat(locale, options)));
  return format;
}

const decimal = new Intl.NumberFormat("hr-HR", { maximumFractionDigits: 1 });
const integer = new Intl.NumberFormat("hr-HR", { maximumFractionDigits: 0 });

/** 14:27 */
export const formatClock = (date: Date, timeZone: string) =>
  dateFormat("hr-HR", { hour: "2-digit", minute: "2-digit", hourCycle: "h23", timeZone }).format(date);

/** 7 h */
export const formatHour = (date: Date, timeZone: string) =>
  `${Number.parseInt(dateFormat("hr-HR", { hour: "numeric", hourCycle: "h23", timeZone }).format(date), 10)} h`;

/** YYYY-MM-DD on the location's calendar. */
export const isoDateOf = (date: Date, timeZone: string) => dateFormat("en-CA", { timeZone }).format(date);

/** Minutes since midnight on the location's clock. */
export function minutesOfDay(date: Date, timeZone: string): number {
  const [h, m] = formatClock(date, timeZone).split(":").map(Number);
  return h * 60 + m;
}

/** "Danas", "Sutra", then the short weekday and date ("Čet 1. 10.") for a YYYY-MM-DD date. */
export function dayLabel(date: string, now: Date, timeZone: string): string {
  const days = Math.round((Date.parse(date) - Date.parse(isoDateOf(now, timeZone))) / 86_400_000);
  if (days === 0) return "Danas";
  if (days === 1) return "Sutra";
  const weekday = dateFormat("hr-HR", { weekday: "short", timeZone: "UTC" })
    .format(new Date(`${date}T12:00:00Z`))
    .replace(".", "");
  const [, month, day] = date.split("-").map(Number); // Intl pads them: "01. 10."
  return `${weekday.charAt(0).toUpperCase()}${weekday.slice(1)} ${day}. ${month}.`;
}

/** 17°, rounded; never "-0°"; a dash when missing. */
export const formatDegrees = (value: number | null | undefined) =>
  value == null ? "–" : `${Math.round(value) || 0}°`;

/** 48 % */
export const formatPercent = (value: number) => `${Math.round(value)} %`;

/** 4,3 mm */
export const formatMillimetres = (value: number) => `${decimal.format(value)} mm`;

/** 12 km/h */
export const formatSpeed = (kmh: number) => `${Math.round(kmh)} km/h`;

/** 1.013 hPa */
export const formatPressure = (hpa: number) => `${integer.format(hpa)} hPa`;

/** 41 km, 2,5 km, 800 m */
export function formatDistance(metres: number): string {
  if (metres < 1000) return `${Math.round(metres / 100) * 100} m`;
  return `${(metres < 10_000 ? decimal : integer).format(metres / 1000)} km`;
}

/** 11 h 33 min */
export function formatDuration(seconds: number): string {
  const minutes = Math.round(seconds / 60);
  return `${Math.floor(minutes / 60)} h ${minutes % 60} min`;
}

/** "upravo", "prije 5 min", "prije 2 h", "prije 3 dana" */
export function formatAgo(then: Date, now: Date): string {
  const minutes = Math.floor((now.getTime() - then.getTime()) / 60_000);
  if (minutes < 1) return "upravo";
  if (minutes < 60) return `prije ${minutes} min`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `prije ${hours} h`;
  const days = Math.floor(hours / 24);
  return `prije ${days} ${days % 10 === 1 && days % 100 !== 11 ? "dan" : "dana"}`;
}
