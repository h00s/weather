<script lang="ts">
  import { oklchCss, type Sky } from "$lib/helpers/sky";
  import type { Precipitation } from "$lib/helpers/weather";
  import PrecipitationLayer from "./precipitation-layer.svelte";

  type Props = {
    sky: Sky;
    precipitation: Precipitation;
    /** How long a color change takes; the sky is recomputed every minute, so a minute makes the
     *  day one continuous fade. 0 in design previews. */
    transitionMs?: number;
  };
  let { sky, precipitation, transitionMs = 60_000 }: Props = $props();
</script>

<div
  class="sky fixed inset-0 -z-10"
  style:--sky-top={oklchCss(sky.top)}
  style:--sky-bottom={oklchCss(sky.bottom)}
  style:transition-duration="{transitionMs}ms"
>
  {#if precipitation !== "none"}
    <PrecipitationLayer kind={precipitation} />
  {/if}
</div>

<style>
  .sky {
    background: linear-gradient(to bottom, var(--sky-top), var(--sky-bottom));
    transition-property: --sky-top, --sky-bottom;
    transition-timing-function: linear;
  }
</style>
