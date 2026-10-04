<script lang="ts">
  import Navigation2 from "@lucide/svelte/icons/navigation-2";
  import {
    formatClock,
    formatDegrees,
    formatDistance,
    formatDuration,
    formatPercent,
    formatPressure,
    formatSpeed,
  } from "$lib/helpers/format";
  import { uvBand } from "$lib/helpers/uv";
  import { compassName, windArrowDegrees, windStrength } from "$lib/helpers/wind";
  import type { DailyForecast, Forecast } from "$lib/types/forecast";
  import DetailTile from "./detail-tile.svelte";

  type Props = { forecast: Forecast; today: DailyForecast | undefined };
  let { forecast, today }: Props = $props();

  const c = $derived(forecast.current);
  const tz = $derived(forecast.timezone);
  const uv = $derived(uvBand(c.uvIndex));
</script>

<section aria-labelledby="details-title">
  <h2 id="details-title" class="sr-only">Pojedinosti</h2>
  <dl class="grid grid-cols-2 gap-3 sm:grid-cols-3">
    <DetailTile
      label="Vjetar"
      value={formatSpeed(c.windSpeed)}
      detail="{windStrength(c.windSpeed)}, {compassName(c.windDirection)} · udari {formatSpeed(c.windGusts)}"
    >
      {#snippet icon()}
        <Navigation2 class="size-5 shrink-0 fill-current" style="transform: rotate({windArrowDegrees(c.windDirection)}deg)" aria-hidden="true" />
      {/snippet}
    </DetailTile>
    <DetailTile
      label="UV indeks"
      value={String(Math.round(c.uvIndex))}
      detail={today ? `${uv.label} · danas do ${Math.round(today.uvIndexMax)}` : uv.label}
    />
    <DetailTile label="Vlažnost" value={formatPercent(c.humidity)} detail="Rosište {formatDegrees(c.dewPoint)}" />
    <DetailTile label="Tlak" value={formatPressure(c.pressure)} detail="Na razini mora" />
    <DetailTile label="Vidljivost" value={formatDistance(c.visibility)} />
    {#if today}
      <DetailTile
        label="Sunce"
        small
        value="{formatClock(new Date(today.sunrise), tz)} – {formatClock(new Date(today.sunset), tz)}"
        detail="Dan traje {formatDuration(today.daylightSeconds)}"
      />
    {/if}
  </dl>
</section>
