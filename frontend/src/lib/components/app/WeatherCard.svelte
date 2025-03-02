<script lang="ts">
  import type { WeatherCard } from "$type/WeatherCard";
  import { formatTimeWithWeekday } from "$util/datetime";
  import { formatTemperature, formatHumidity, formatPrecipitation, weatherIcon } from "$util/weather";

  const { weather }: { weather: WeatherCard } = $props<{
    weather: WeatherCard;
  }>();

  const formattedTime = $derived(formatTimeWithWeekday(weather.timestamp));
  const temperature: string = $derived(formatTemperature(weather.temperature));
  const humidity: string = $derived(formatHumidity(weather.humidity));
  const precipitation: string = $derived(formatPrecipitation(weather.precipitation));
</script>

<div class="grid grid-cols-3 items-center relative">
  <!-- Left section - City and time -->
  <div class="flex flex-col">
    <h2 class="text-3xl font-semibold text-white">{weather.city}</h2>
    <p class="text-blue-300 text-sm mt-1">{formattedTime}</p>
  </div>

  <!-- Middle section - Weather icon -->
  <div class="flex justify-center text-5xl">
    {weatherIcon(weather.code)}
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
