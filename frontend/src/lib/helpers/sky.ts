import type { Precipitation } from "./weather";

export interface Oklch {
  l: number;
  c: number;
  h: number;
}

export interface Sky {
  top: Oklch;
  bottom: Oklch;
}

const MIN = 60_000;

const oklch = (l: number, c: number, h: number): Oklch => ({ l, c, h });

// The sky over a day in the brand's navy-to-teal family: [anchor, minutes from it, zenith, horizon].
// Night settles on the brand navy (#012a4a), the day sits between #01497c and #2c7da0, and only the
// hour around sunrise and sunset warms the horizon. skyAt darkens anything that would leave white
// text below MIN_CONTRAST.
const NIGHT: Sky = { top: oklch(0.18, 0.05, 252), bottom: oklch(0.278, 0.073, 248) };
const DAY: Sky = { top: oklch(0.4, 0.106, 248), bottom: oklch(0.557, 0.094, 231) };
const keyframes: ["sunrise" | "sunset", number, Sky][] = [
  ["sunrise", -70, NIGHT],
  ["sunrise", -25, { top: oklch(0.26, 0.07, 270), bottom: oklch(0.46, 0.1, 15) }], // indigo over a muted rose
  ["sunrise", 15, { top: oklch(0.35, 0.09, 255), bottom: oklch(0.54, 0.1, 55) }], // navy over apricot
  ["sunrise", 75, DAY],
  ["sunset", -100, DAY],
  ["sunset", -35, { top: oklch(0.38, 0.1, 252), bottom: oklch(0.55, 0.1, 65) }], // golden hour
  ["sunset", 5, { top: oklch(0.3, 0.085, 265), bottom: oklch(0.47, 0.12, 35) }], // afterglow
  ["sunset", 35, { top: oklch(0.21, 0.06, 258), bottom: oklch(0.31, 0.07, 270) }], // dusk
  ["sunset", 75, NIGHT],
];

/** The contrast white text keeps against any sky. */
export const MIN_CONTRAST = 4.5;

/** The sky gradient for now: interpolated between the day's keyframes, then greyed by cloud cover
 *  (percent) and darkened by precipitation. */
export function skyAt(now: Date, sunrise: Date, sunset: Date, cloudCover: number, precipitation: Precipitation): Sky {
  const t = now.getTime();
  const frames = keyframes.map(([anchor, minutes, sky]) => ({
    t: (anchor === "sunrise" ? sunrise : sunset).getTime() + minutes * MIN,
    sky,
  }));

  let sky = NIGHT;
  for (let i = 0; i < frames.length - 1; i++) {
    const a = frames[i];
    const b = frames[i + 1];
    if (t >= a.t && t < b.t) {
      const f = (t - a.t) / (b.t - a.t);
      sky = { top: mix(a.sky.top, b.sky.top, f), bottom: mix(a.sky.bottom, b.sky.bottom, f) };
      break;
    }
  }

  const cloud = Math.min(Math.max(cloudCover, 0), 100) / 100;
  const dim = { none: 1, snow: 0.95, drizzle: 0.9, rain: 0.85, storm: 0.75 }[precipitation];
  const weather = (c: Oklch): Oklch => ({ l: c.l * (1 - 0.1 * cloud) * dim, c: c.c * (1 - 0.75 * cloud), h: c.h });
  return { top: legible(weather(sky.top)), bottom: legible(weather(sky.bottom)) };
}

/** Darkens c just enough that white text on it keeps MIN_CONTRAST. */
function legible(c: Oklch): Oklch {
  let l = c.l;
  while (l > 0 && contrastWithWhite({ ...c, l }) < MIN_CONTRAST) l -= 0.005;
  return { ...c, l };
}

/** Interpolates in OKLab, so blue to orange passes through a pale tone rather than green. */
function mix(a: Oklch, b: Oklch, f: number): Oklch {
  const rad = Math.PI / 180;
  const [aa, ab] = [a.c * Math.cos(a.h * rad), a.c * Math.sin(a.h * rad)];
  const [ba, bb] = [b.c * Math.cos(b.h * rad), b.c * Math.sin(b.h * rad)];
  const x = aa + (ba - aa) * f;
  const y = ab + (bb - ab) * f;
  const h = (Math.atan2(y, x) / rad + 360) % 360;
  return { l: a.l + (b.l - a.l) * f, c: Math.hypot(x, y), h };
}

export const oklchCss = ({ l, c, h }: Oklch) => `oklch(${l.toFixed(3)} ${c.toFixed(3)} ${h.toFixed(1)})`;

/** Linear-light sRGB of c, each channel clamped to 0–1 (OKLab → sRGB, Björn Ottosson's matrices). */
function linearSrgb(c: Oklch): [number, number, number] {
  const rad = Math.PI / 180;
  const a = c.c * Math.cos(c.h * rad);
  const b = c.c * Math.sin(c.h * rad);
  const l = (c.l + 0.3963377774 * a + 0.2158037573 * b) ** 3;
  const m = (c.l - 0.1055613458 * a - 0.0638541728 * b) ** 3;
  const s = (c.l - 0.0894841775 * a - 1.291485548 * b) ** 3;
  const clamp = (v: number) => Math.min(1, Math.max(0, v));
  return [
    clamp(4.0767416621 * l - 3.3077115913 * m + 0.2309699292 * s),
    clamp(-1.2684380046 * l + 2.6097574011 * m - 0.3413193965 * s),
    clamp(-0.0041960863 * l - 0.7034186147 * m + 1.707614701 * s),
  ];
}

/** WCAG contrast ratio of white text on c (1–21). */
export function contrastWithWhite(c: Oklch): number {
  const [r, g, b] = linearSrgb(c);
  return 1.05 / (0.2126 * r + 0.7152 * g + 0.0722 * b + 0.05);
}

/** c as #rrggbb, for what takes no oklch(): the theme-color meta. */
export function oklchToHex(c: Oklch): string {
  return `#${linearSrgb(c)
    .map((v) => (v <= 0.0031308 ? 12.92 * v : 1.055 * v ** (1 / 2.4) - 0.055))
    .map((v) => Math.round(Math.min(1, Math.max(0, v)) * 255).toString(16).padStart(2, "0"))
    .join("")}`;
}
