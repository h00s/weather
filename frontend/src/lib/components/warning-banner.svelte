<script lang="ts">
  import TriangleAlert from "@lucide/svelte/icons/triangle-alert";
  import { warningKey, warningWhen } from "$lib/helpers/alerts";
  import type { Warning } from "$lib/types/warning";

  type Props = { warnings: Warning[]; now: Date; timeZone: string };
  let { warnings, now, timeZone }: Props = $props();

  // Status colours, always with the icon and DHMZ's own words ("Žuto upozorenje za vjetar").
  const fill = { yellow: "bg-warn-yellow text-black", orange: "bg-warn-orange text-black", red: "bg-warn-red text-white" };
</script>

<section aria-label="Upozorenja DHMZ-a" class="flex flex-col gap-2">
  {#each warnings as w (warningKey(w))}
    <details class="rounded-2xl {fill[w.level]}">
      <summary class="flex cursor-pointer list-none items-center gap-3 px-4 py-2.5 font-medium [&::-webkit-details-marker]:hidden">
        <TriangleAlert class="size-5 shrink-0" aria-hidden="true" />
        <span class="min-w-0 flex-1">{w.event}</span>
        <span class="shrink-0 text-sm opacity-80">{warningWhen(w, now, timeZone)}</span>
      </summary>
      {#if w.description}<p class="px-4 pb-3 text-sm">{w.description}</p>{/if}
    </details>
  {/each}
</section>
