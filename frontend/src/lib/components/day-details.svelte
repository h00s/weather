<script lang="ts">
  import {
    formatClock,
    formatDegrees,
    formatDuration,
    formatHour,
    formatMillimetres,
    formatPercent,
    formatSpeed,
    minutesOfDay,
  } from "$lib/helpers/format";
  import { uvBand } from "$lib/helpers/uv";
  import { weatherInfo } from "$lib/helpers/weather";
  import { compassName } from "$lib/helpers/wind";
  import type { DailyForecast, HourlyForecast } from "$lib/types/forecast";
  import ScrollStrip from "./scroll-strip.svelte";
  import WeatherIcon from "./weather-icon.svelte";

  type Props = { day: DailyForecast; hours: HourlyForecast[]; timeZone: string };
  let { day, hours, timeZone }: Props = $props();

  // Every third hour keeps the strip short enough for a phone.
  const shown = $derived(hours.filter((h) => (minutesOfDay(new Date(h.time), timeZone) / 60) % 3 === 0));
</script>

<div class="px-2 pt-1 pb-3">
  {#if shown.length}
    <ScrollStrip label="Prognoza po satima za taj dan">
      <ol class="flex w-max gap-1 pb-1">
        {#each shown as h (h.time)}
          <li class="flex min-w-14 flex-col items-center rounded-xl bg-white/6 px-1 py-2 text-sm">
            <span class="text-muted-foreground">{formatHour(new Date(h.time), timeZone)}</span>
            <WeatherIcon name={weatherInfo(h.weatherCode, h.isDay).icon} size={32} />
            <span class="font-medium">{formatDegrees(h.temperature)}</span>
            {#if h.precipitationProbability >= 20}<span class="text-rain text-xs">{formatPercent(h.precipitationProbability)}</span>{/if}
          </li>
        {/each}
      </ol>
    </ScrollStrip>
  {/if}
  <dl class="mt-2 grid grid-cols-1 gap-x-4 gap-y-1.5 text-sm sm:grid-cols-2">
    <div><dt class="text-muted-foreground inline">Oborine</dt> <dd class="inline">{formatMillimetres(day.precipitationSum)} · {formatPercent(day.precipitationProbabilityMax)}</dd></div>
    <div><dt class="text-muted-foreground inline">Vjetar</dt> <dd class="inline">do {formatSpeed(day.windSpeedMax)}, {compassName(day.windDirectionDominant)} · udari {formatSpeed(day.windGustsMax)}</dd></div>
    <div><dt class="text-muted-foreground inline">UV indeks</dt> <dd class="inline">{Math.round(day.uvIndexMax)} · {uvBand(day.uvIndexMax).label}</dd></div>
    <div><dt class="text-muted-foreground inline">Sunce</dt> <dd class="inline">{formatClock(new Date(day.sunrise), timeZone)} – {formatClock(new Date(day.sunset), timeZone)} ({formatDuration(day.daylightSeconds)})</dd></div>
  </dl>
</div>
