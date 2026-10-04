export type Series = (number | null)[];

/** The lowest and highest value, or null when there is none. */
export function extent(values: Series): [number, number] | null {
  let lo = Infinity;
  let hi = -Infinity;
  for (const v of values) {
    if (v == null) continue;
    lo = Math.min(lo, v);
    hi = Math.max(hi, v);
  }
  return lo <= hi ? [lo, hi] : null;
}

/** A linear map from domain onto range; a flat domain maps to the middle of the range. */
export function scale([d0, d1]: [number, number], [r0, r1]: [number, number]) {
  return (v: number) => (d0 === d1 ? (r0 + r1) / 2 : r0 + ((v - d0) / (d1 - d0)) * (r1 - r0));
}

const num = (n: number) => String(Math.round(n * 100) / 100);

/** Points of each unbroken run of values, evenly spread across width. pad keeps the extremes off
 *  the edges, so a 2px stroke isn't clipped. */
function runs(values: Series, width: number, height: number, domain: [number, number], pad: number) {
  const x = (i: number) => (values.length < 2 ? width / 2 : (i * width) / (values.length - 1));
  const y = scale(domain, [height - pad, pad]);
  const out: [number, number][][] = [];
  let run: [number, number][] = [];
  values.forEach((v, i) => {
    if (v == null) {
      if (run.length) out.push(run);
      run = [];
    } else {
      run.push([x(i), y(v)]);
    }
  });
  if (run.length) out.push(run);
  return out;
}

/** An SVG path through the values; a null breaks the line. */
export function linePath(values: Series, width: number, height: number, domain: [number, number], pad = 0): string {
  return runs(values, width, height, domain, pad)
    .map((run) => run.map(([px, py], i) => `${i ? "L" : "M"}${num(px)},${num(py)}`).join(""))
    .join("");
}

/** The area under each run of two or more values, closed down to the bottom edge. */
export function areaPath(values: Series, width: number, height: number, domain: [number, number], pad = 0): string {
  return runs(values, width, height, domain, pad)
    .filter((run) => run.length > 1)
    .map((run) => {
      const line = run.map(([px, py], i) => `${i ? "L" : "M"}${num(px)},${num(py)}`).join("");
      return `${line}L${num(run[run.length - 1][0])},${num(height)}L${num(run[0][0])},${num(height)}Z`;
    })
    .join("");
}
