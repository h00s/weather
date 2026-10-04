<script lang="ts">
  import RotateCw from "@lucide/svelte/icons/rotate-cw";
  import type { Snippet } from "svelte";

  type Props = { onrefresh: () => Promise<void>; children: Snippet };
  let { onrefresh, children }: Props = $props();

  // Pulling THRESHOLD px (after 1:2 resistance) at the top of the page and letting go refreshes.
  const THRESHOLD = 72;
  const MAX = 112;

  let pull = $state(0);
  let dragging = $state(false);
  let refreshing = $state(false);
  let startY = 0;

  function ontouchstart(e: TouchEvent) {
    if (refreshing || window.scrollY > 0 || e.touches.length !== 1) return;
    startY = e.touches[0].clientY;
    dragging = true;
  }

  function ontouchmove(e: TouchEvent) {
    if (!dragging) return;
    const dy = e.touches[0].clientY - startY;
    pull = dy > 0 && window.scrollY <= 0 ? Math.min(MAX, dy / 2) : 0;
  }

  async function ontouchend() {
    if (!dragging) return;
    dragging = false;
    if (pull < THRESHOLD) {
      pull = 0;
      return;
    }
    refreshing = true;
    pull = THRESHOLD;
    try {
      await onrefresh();
    } finally {
      refreshing = false;
      pull = 0;
    }
  }
</script>

<svelte:window {ontouchstart} {ontouchmove} {ontouchend} ontouchcancel={ontouchend} />

<!-- Touch only, and decorative: the header's refresh button does the same for everyone. -->
<div
  aria-hidden="true"
  class="pointer-events-none fixed inset-x-0 top-0 z-20 flex justify-center"
  style:transform="translateY({pull - 44}px)"
  style:opacity={Math.min(1, pull / THRESHOLD)}
  style:transition={dragging ? "none" : "transform 200ms ease, opacity 200ms ease"}
>
  <div class="mt-[max(0.5rem,env(safe-area-inset-top))] grid size-10 place-items-center rounded-full bg-white/20 shadow-lg backdrop-blur">
    <RotateCw
      class="size-5 {refreshing ? 'animate-spin' : ''}"
      style={refreshing ? undefined : `transform: rotate(${(pull / THRESHOLD) * 270}deg)`}
    />
  </div>
</div>

<!-- Transformed only while pulled: a transform would otherwise trap fixed-position descendants. -->
<div style:transform={pull ? `translateY(${pull / 2}px)` : undefined} style:transition={dragging ? "none" : "transform 200ms ease"}>
  {@render children()}
</div>
