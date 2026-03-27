<script lang="ts">
  import type { CurrentForecast } from "$type/Forecast";
  import { formatTimeWithWeekday } from "$util/datetime";
  import { formatTemperature, formatHumidity, formatPrecipitation, weatherIcon } from "$util/weather";

  const { city, data }: {
    city: string;
    data: CurrentForecast;
  } = $props<{
    city: string;
    data: CurrentForecast;
  }>();

  const formattedTime: string = $derived(formatTimeWithWeekday(data.time));
  const temperature: string = $derived(formatTemperature(data.data.temperature_2m));
  const humidity: string = $derived(formatHumidity(data.data.relative_humidity_2m));
  const precipitation: string = $derived(formatPrecipitation(data.data.precipitation));
</script>

<div class="flex flex-col items-center text-center py-2 md:py-6">
  <p class="text-blue-300 text-xs md:text-sm uppercase tracking-widest font-medium">{city}</p>
  <p class="text-blue-300/60 text-xs md:text-sm mt-1">{formattedTime}</p>

  <div class="text-7xl md:text-8xl lg:text-9xl my-5 md:my-7 leading-none">
    {weatherIcon(data.data.weather_code)}
  </div>

  <div class="text-6xl md:text-7xl lg:text-8xl font-thin text-white tracking-tight">
    {temperature}
  </div>

  <div class="flex items-center gap-3 mt-4 md:mt-6 text-sm md:text-base text-blue-300/80">
    <span>{humidity}</span>
    <span class="text-blue-300/30">·</span>
    <span>{precipitation}</span>
  </div>
</div>
