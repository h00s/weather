const points = ["S", "SI", "I", "JI", "J", "JZ", "Z", "SZ"];
const names = ["sjever", "sjeveroistok", "istok", "jugoistok", "jug", "jugozapad", "zapad", "sjeverozapad"];
const index = (degrees: number) => Math.round((((degrees % 360) + 360) % 360) / 45) % 8;

/** The 8-point compass point the wind blows from, in Croatian: 0° "S", 45° "SI", 225° "JZ". */
export const compassPoint = (degrees: number) => points[index(degrees)];

/** The same, spelled out: "jugozapad". */
export const compassName = (degrees: number) => names[index(degrees)];

/** Rotation for an arrow drawn pointing up, so it points where the wind blows to. */
export const windArrowDegrees = (degrees: number) => (degrees + 180) % 360;

/** The Beaufort groups in Croatian, by km/h. */
export function windStrength(kmh: number): string {
  if (kmh < 2) return "Tiho";
  if (kmh < 20) return "Slab vjetar";
  if (kmh < 39) return "Umjeren vjetar";
  if (kmh < 62) return "Jak vjetar";
  if (kmh < 89) return "Olujni vjetar";
  return "Orkanski vjetar";
}
