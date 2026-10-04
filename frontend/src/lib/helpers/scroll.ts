/** Whether a horizontal strip has more to show before (start) and after (end) its view. */
export interface ScrollEdges {
  start: boolean;
  end: boolean;
}

/** Sub-pixel scroll positions (zoom, fractional widths) never count as "more". */
const SLACK = 1;

export function scrollEdges(scrollLeft: number, clientWidth: number, scrollWidth: number): ScrollEdges {
  return {
    start: scrollLeft > SLACK,
    end: scrollLeft + clientWidth < scrollWidth - SLACK,
  };
}

/** How far one arrow press moves: most of a view, so the last column stays for context. */
export const scrollStep = (clientWidth: number) => Math.round(clientWidth * 0.8);

/** A CSS mask that fades out the sides with more content: the hint to swipe or scroll. */
export function fadeMask({ start, end }: ScrollEdges): string {
  if (!start && !end) return "none";
  const from = start ? "transparent 0, #000 3rem" : "#000 0";
  const to = end ? "#000 calc(100% - 3rem), transparent 100%" : "#000 100%";
  return `linear-gradient(to right, ${from}, ${to})`;
}
