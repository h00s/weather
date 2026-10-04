<script lang="ts">
  import { activePollen, aqiLabels, pollenLevelLabels, pollenNames } from "$lib/helpers/air";
  import type { AirQuality } from "$lib/types/airquality";

  type Props = { airQuality: AirQuality };
  let { airQuality }: Props = $props();

  const pollen = $derived(activePollen(airQuality));
  // The European AQI's own colour scale, always next to the label that says the same.
  const aqiDot = {
    good: "bg-teal-300",
    fair: "bg-emerald-300",
    moderate: "bg-yellow-300",
    poor: "bg-orange-400",
    very_poor: "bg-red-500",
    extremely_poor: "bg-purple-500",
  };
  const pollenDot = { none: "", low: "bg-emerald-300", moderate: "bg-yellow-300", high: "bg-orange-400", very_high: "bg-red-500" };
</script>

<section aria-labelledby="air-title" class="panel p-5">
  <h2 id="air-title" class="text-muted-foreground text-xs font-medium tracking-wider uppercase">Kvaliteta zraka</h2>
  <p class="mt-2 flex items-center gap-3">
    <span class="size-3 shrink-0 rounded-full {aqiDot[airQuality.aqi.level]}" aria-hidden="true"></span>
    <span class="text-2xl font-light">{aqiLabels[airQuality.aqi.level]}</span>
    <span class="text-muted-foreground">EAQI {Math.round(airQuality.aqi.value)}</span>
  </p>
  <h3 class="text-muted-foreground mt-4 text-xs font-medium tracking-wider uppercase">Pelud danas</h3>
  {#if pollen.length}
    <ul class="mt-2 flex flex-wrap gap-2">
      {#each pollen as p (p.type)}
        <li class="flex items-center gap-2 rounded-full bg-white/10 px-3 py-1 text-sm">
          <span class="size-2 rounded-full {pollenDot[p.level]}" aria-hidden="true"></span>
          {pollenNames[p.type]} <span class="text-muted-foreground">{pollenLevelLabels[p.level]}</span>
        </li>
      {/each}
    </ul>
  {:else}
    <p class="text-muted-foreground mt-1">Nema peludi u zraku.</p>
  {/if}
</section>
