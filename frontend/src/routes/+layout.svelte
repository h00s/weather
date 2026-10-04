<script lang="ts">
  import "./layout.css";
  import { updated } from "$app/state";
  import type { LayoutProps } from "./$types";

  let { children }: LayoutProps = $props();

  // SvelteKit applies a new deploy on the next navigation, and this one-page app rarely navigates.
  // When the version poll (vite.config.ts) sees a new build, reload the next time the page comes
  // back into view: never under the reader's eyes. A full reload is the point here, hence location.
  $effect(() => {
    if (!updated.current) return;
    const reload = () => {
      if (document.visibilityState === "visible") location.reload();
    };
    document.addEventListener("visibilitychange", reload);
    return () => document.removeEventListener("visibilitychange", reload);
  });
</script>

{@render children()}
