<script lang="ts">
  import type { DailyForecast } from "$type/Forecast";
  import { formatTemperature, formatPrecipitation, weatherIcon } from "$util/weather";
  import { formatDateWithWeekday } from "$util/datetime";

  const { data }: {
    data: DailyForecast;
  } = $props<{
    data: DailyForecast;
  }>();
</script>

{#if data !== undefined && data.time.length > 0}
  <div>
    <p class="text-blue-300/60 text-xs uppercase tracking-widest mb-3 md:mb-4">7-day forecast</p>
    <div class="divide-y divide-white/[0.06]">
      {#each data.time as _, index}
        <div class="flex items-center gap-3 py-2.5 md:py-3">
          <span class="text-white/80 text-sm md:text-base w-24 md:w-28 shrink-0">{formatDateWithWeekday(data.time[index])}</span>
          <span class="text-2xl md:text-3xl leading-none">{weatherIcon(data.data.weather_code[index])}</span>
          <span class="text-blue-300/60 text-xs md:text-sm ml-auto shrink-0">{formatPrecipitation(data.data.precipitation_sum[index], true)}</span>
          <span class="text-blue-300/70 text-sm md:text-base shrink-0">{formatTemperature(data.data.temperature_2m_min[index])}</span>
          <span class="text-white text-sm md:text-base font-medium shrink-0">{formatTemperature(data.data.temperature_2m_max[index])}</span>
        </div>
      {/each}
    </div>
  </div>
{/if}
