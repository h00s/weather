<script lang="ts">
  import type { CurrentForecast } from "$type/Forecast";
  import { formatTimeWithWeekday } from "$util/datetime";
  import { formatTemperature, formatHumidity, formatPrecipitation, weatherIcon } from "$util/weather";

  const { city, data }: { city: string, data: CurrentForecast } = $props<{
    city: string;
    data: CurrentForecast;
  }>();

  const formattedTime = $derived(formatTimeWithWeekday(data.time));
  const temperature: string = $derived(formatTemperature(data.data.temperature_2m));
  const humidity: string = $derived(formatHumidity(data.data.relative_humidity_2m));
  const precipitation: string = $derived(formatPrecipitation(data.data.precipitation));
</script>

<div class="grid grid-cols-3 items-center relative">
  <!-- Left section - City and time -->
  <div class="flex flex-col">
    <h2 class="text-3xl font-semibold text-white">{city}</h2>
    <p class="text-blue-300 text-sm mt-1">{formattedTime}</p>
  </div>

  <!-- Middle section - Weather icon -->
  <div class="flex justify-center text-5xl">
    {weatherIcon(data.data.weather_code)}
  </div>

  <!-- Right section - Temperature and conditions -->
  <div class="flex flex-col items-end">
    <div class="text-3xl font-bold text-white">
      {temperature}
    </div>
    <p class="text-blue-300 text-sm">
      {humidity},
      {precipitation}
    </p>
  </div>
</div>
