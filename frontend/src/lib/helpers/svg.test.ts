import { describe, expect, it } from "vitest";
import { stripAnimation } from "./svg";

describe("stripAnimation", () => {
  it("removes SMIL elements, self-closing or not, and keeps the drawing", () => {
    const svg =
      '<svg><g><path d="M0 0"/><animateTransform attributeName="transform" type="rotate" values="0;45" dur="6s" repeatCount="indefinite"/></g>' +
      '<circle r="2"><animate attributeName="r" values="2;3" dur="1s"></animate></circle><set attributeName="opacity" to="1"/></svg>';
    expect(stripAnimation(svg)).toBe('<svg><g><path d="M0 0"/></g><circle r="2"></circle></svg>');
  });
});
