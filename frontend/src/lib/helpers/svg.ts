/** The SVG without its SMIL animation (<animate>, <animateTransform>, <animateMotion>, <set>). An
 *  <img> plays SMIL whatever CSS says, so readers who asked for reduced motion get this still. */
export function stripAnimation(svg: string): string {
  return svg
    .replace(/<(animate|animateTransform|animateMotion|set)\b[^>]*\/>/g, "")
    .replace(/<(animate|animateTransform|animateMotion|set)\b[^>]*>[\s\S]*?<\/\1>/g, "");
}

const stills = new Map<string, Promise<string>>();

/** A data: URL of the still version of the SVG at url, fetched once per url. */
export function stillSvgUrl(url: string): Promise<string> {
  let still = stills.get(url);
  if (!still) {
    still = fetch(url)
      .then((r) => r.text())
      .then((svg) => `data:image/svg+xml,${encodeURIComponent(stripAnimation(svg))}`);
    stills.set(url, still);
    still.catch(() => stills.delete(url)); // retry next time
  }
  return still;
}
