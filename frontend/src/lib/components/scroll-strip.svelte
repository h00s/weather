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
<div class={["relative", className]}>
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
    <button type="button" tabindex="-1" aria-hidden="true" onclick={() => page(-1)} class="strip-arrow left-1">
      <ChevronLeft class="size-5" />
    </button>
  {/if}
  {#if edges.end}
    <button type="button" tabindex="-1" aria-hidden="true" onclick={() => page(1)} class="strip-arrow right-1">
      <ChevronRight class="size-5" />
    </button>
  {/if}
</div>

<style>
  /* Buttons only where there is a mouse to click them; a touch screen swipes. */
  .strip-arrow {
    display: none;
  }
  @media (hover: hover) and (pointer: fine) {
    .strip-arrow {
      position: absolute;
      top: 50%;
      translate: 0 -50%;
      display: grid;
      place-items: center;
      width: 2.25rem;
      height: 2.25rem;
      border-radius: 9999px;
      color: var(--color-foreground);
      background-color: oklch(0.25 0.06 248 / 85%);
      box-shadow:
        inset 0 0 0 1px oklch(1 0 0 / 18%),
        0 4px 14px oklch(0 0 0 / 30%);
      -webkit-backdrop-filter: blur(6px);
      backdrop-filter: blur(6px);
      cursor: pointer;
      transition: background-color 150ms;
    }
    .strip-arrow:hover {
      background-color: oklch(0.32 0.07 248 / 95%);
    }
  }
</style>
