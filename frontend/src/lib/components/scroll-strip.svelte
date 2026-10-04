<script lang="ts">
  import ChevronLeft from "@lucide/svelte/icons/chevron-left";
  import ChevronRight from "@lucide/svelte/icons/chevron-right";
  import type { Snippet } from "svelte";
  import { prefersReducedMotion } from "svelte/motion";
  import { fadeMask, scrollEdges, scrollStep, type ScrollEdges } from "$lib/helpers/scroll";

  type Props = {
    /** Names the strip for screen readers and keyboard users, who scroll it with the arrow keys. */
    label: string;
    children: Snippet;
    class?: string;
  };
  let { label, children, class: className }: Props = $props();

  let viewport = $state<HTMLDivElement>();
  let edges = $state<ScrollEdges>({ start: false, end: false });

  function measure() {
    if (viewport) edges = scrollEdges(viewport.scrollLeft, viewport.clientWidth, viewport.scrollWidth);
  }

  // The strip and its content change size with the window and with new forecasts.
  $effect(() => {
    if (!viewport) return;
    const observer = new ResizeObserver(measure);
    observer.observe(viewport);
    if (viewport.firstElementChild) observer.observe(viewport.firstElementChild);
    measure();
    return () => observer.disconnect();
  });

  function page(direction: -1 | 1) {
    if (!viewport) return;
    viewport.scrollBy({
      left: direction * scrollStep(viewport.clientWidth),
      behavior: prefersReducedMotion.current ? "auto" : "smooth",
    });
  }
</script>

<!-- No scrollbar: the faded edge says there is more. A touch screen swipes, a trackpad or the arrow
     keys scroll, and with a mouse the round buttons page through. -->
<div class={["strip relative", className]}>
  <!-- svelte-ignore a11y_no_noninteractive_tabindex -->
  <div
    bind:this={viewport}
    onscroll={measure}
    tabindex="0"
    role="region"
    aria-label={label}
    class="scrollbar-none overflow-x-auto"
    style:mask-image={fadeMask(edges)}
    style:-webkit-mask-image={fadeMask(edges)}
  >
    {@render children()}
  </div>
  {#if edges.start}
    <button type="button" tabindex="-1" aria-hidden="true" onclick={() => page(-1)} class="strip-arrow left-3">
      <ChevronLeft class="size-4" />
    </button>
  {/if}
  {#if edges.end}
    <button type="button" tabindex="-1" aria-hidden="true" onclick={() => page(1)} class="strip-arrow right-3">
      <ChevronRight class="size-4" />
    </button>
  {/if}
</div>

<style>
  /* Buttons only where there is a mouse to click them; a touch screen swipes. */
  .strip-arrow {
    display: none;
  }
  /* Frosted glass like the card, quiet until the strip is hovered: a hint, not a call to action. */
  @media (hover: hover) and (pointer: fine) {
    .strip-arrow {
      position: absolute;
      top: 50%;
      translate: 0 -50%;
      display: grid;
      place-items: center;
      width: 2rem;
      height: 2rem;
      border-radius: 9999px;
      color: oklch(1 0 0 / 90%);
      background-color: oklch(1 0 0 / 12%);
      box-shadow: inset 0 0 0 1px oklch(1 0 0 / 22%);
      -webkit-backdrop-filter: blur(8px);
      backdrop-filter: blur(8px);
      opacity: 0.7;
      cursor: pointer;
      transition:
        opacity 150ms,
        background-color 150ms;
    }
    .strip:hover .strip-arrow {
      opacity: 1;
    }
    .strip-arrow:hover {
      background-color: oklch(1 0 0 / 22%);
    }
  }
</style>
