<script lang="ts">
  import type { HourlyForecast } from "$type/Forecast";
  import { formatTime } from "$util/datetime";
  import {
    formatTemperature,
    formatHumidity,
    weatherIcon,
  } from "$util/weather";

  const {
    data,
  }: {
    data: HourlyForecast;
  } = $props<{
    data: HourlyForecast;
  }>();

  const forecast = $derived(processHourlyForecast(data));

  function processHourlyForecast(data: HourlyForecast): HourlyForecast {
    const now = new Date();
    const currentHour = now.getHours();

    const start = data.time.findIndex((timeStr) => {
      const timeDate = new Date(timeStr);
      return timeDate.getHours() >= currentHour;
    });
    const end = start + 12;

    return {
      time: data.time.slice(start, end),
      data: {
        precipitation: data.data.precipitation.slice(start, end),
        relative_humidity_2m: data.data.relative_humidity_2m.slice(start, end),
        temperature_2m: data.data.temperature_2m.slice(start, end),
        weather_code: data.data.weather_code.slice(start, end),
      },
    };
  }
</script>

{#if forecast !== undefined && forecast.time.length > 0}
  <div>
    <p class="text-blue-300/60 text-xs uppercase tracking-widest mb-3 md:mb-4">Next 12 hours</p>
    <div class="flex gap-2 overflow-x-auto scrollbar-hidden pb-1">
      {#each forecast.time as _, index}
        <div class="hour-card flex flex-col items-center min-w-[68px] md:min-w-20 py-3 md:py-4 px-2 md:px-3">
          <span class="text-blue-300/70 text-xs md:text-sm">{formatTime(forecast.time[index])}</span>
          <span class="text-2xl md:text-3xl my-2 leading-none">{weatherIcon(forecast.data.weather_code[index])}</span>
          <span class="text-white text-sm md:text-base font-medium">{formatTemperature(forecast.data.temperature_2m[index])}</span>
          <span class="text-blue-300/60 text-xs md:text-sm mt-1">{formatHumidity(forecast.data.relative_humidity_2m[index])}</span>
        </div>
      {/each}
    </div>
  </div>
{/if}
