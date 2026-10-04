/** The WHO band of a UV index, in Croatian, with a level from 0 (low) to 4 (extreme). */
export function uvBand(index: number): { label: string; level: 0 | 1 | 2 | 3 | 4 } {
  const uv = Math.round(index);
  if (uv < 3) return { label: "Nizak", level: 0 };
  if (uv < 6) return { label: "Umjeren", level: 1 };
  if (uv < 8) return { label: "Visok", level: 2 };
  if (uv < 11) return { label: "Vrlo visok", level: 3 };
  return { label: "Ekstreman", level: 4 };
}
