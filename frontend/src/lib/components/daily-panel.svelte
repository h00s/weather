<script lang="ts">
  import ChevronDown from "@lucide/svelte/icons/chevron-down";
  import { hoursOfDay, weekRange } from "$lib/helpers/forecast";
  import { dayLabel, formatDegrees, formatPercent, formatShortDate } from "$lib/helpers/format";
  import { weatherInfo } from "$lib/helpers/weather";
  import type { Forecast } from "$lib/types/forecast";
  import DayDetails from "./day-details.svelte";
  import WeatherIcon from "./weather-icon.svelte";

  type Props = { forecast: Forecast; now: Date };
  let { forecast, now }: Props = $props();

  // Every range bar sits on the week's scale, so the days compare at a glance.
  const range = $derived(weekRange(forecast.daily));
  const pct = (t: number) => Math.min(100, Math.max(0, ((t - range[0]) / Math.max(1, range[1] - range[0])) * 100));
</script>

<section aria-labelledby="daily-title" class="panel px-3 py-4 sm:px-4">
  <h2 id="daily-title" class="text-muted-foreground mb-1 px-2 text-xs font-medium tracking-wider uppercase">7 dana</h2>
  <ul>
    {#each forecast.daily as day (day.date)}
      {@const label = dayLabel(day.date, now, forecast.timezone)}
      {@const info = weatherInfo(day.weatherCode)}
      <li class="border-t border-white/8 first:border-t-0">
        <details class="group">
          <summary class="flex cursor-pointer list-none items-center gap-2 rounded-xl px-2 py-2 hover:bg-white/8 sm:gap-3 [&::-webkit-details-marker]:hidden">
            <span class="w-14 shrink-0 leading-tight">
              <span class="block font-medium">{label.split(" ")[0]}</span>
              <span class="text-muted-foreground block text-xs">{formatShortDate(day.date)}</span>
            </span>
            <!-- The chance of rain sits under the icon, so the range bar keeps its width on a phone. -->
            <span class="flex w-10 shrink-0 flex-col items-center">
              <WeatherIcon name={info.icon} size={32} />
              {#if day.precipitationProbabilityMax >= 20}
                <span class="text-rain text-[11px] leading-none whitespace-nowrap">{formatPercent(day.precipitationProbabilityMax)}</span>
              {/if}
            </span>
            <span class="sr-only">{info.label},</span>
            <span class="text-muted-foreground w-8 shrink-0 text-right tabular-nums"><span class="sr-only">najniža</span>{formatDegrees(day.temperatureMin)}</span>
            <span class="relative mx-1 h-1.5 min-w-12 flex-1 rounded-full bg-white/12" aria-hidden="true">
              <span
                class="absolute inset-y-0 rounded-full bg-linear-to-r from-sky-300 to-amber-300"
                style:left="{pct(day.temperatureMin)}%"
                style:right="{100 - pct(day.temperatureMax)}%"
              ></span>
              {#if label === "Danas"}
                <span class="absolute top-1/2 size-2.5 -translate-x-1/2 -translate-y-1/2 rounded-full bg-white ring-2 ring-black/20" style:left="{pct(forecast.current.temperature)}%"></span>
              {/if}
            </span>
            <span class="w-8 shrink-0 text-right font-medium tabular-nums"><span class="sr-only">najviša</span>{formatDegrees(day.temperatureMax)}</span>
            <ChevronDown class="size-4 shrink-0 opacity-60 transition group-open:rotate-180" aria-hidden="true" />
          </summary>
          <DayDetails {day} hours={hoursOfDay(forecast, day.date)} timeZone={forecast.timezone} />
        </details>
      </li>
    {/each}
  </ul>
</section>
