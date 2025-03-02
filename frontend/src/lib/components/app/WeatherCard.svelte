<script lang="ts">
  import WeatherIcon from "$comp/app/WeatherIcon.svelte";
  import type { WeatherCard } from "$type/WeatherCard";
  import { formatTime } from "$util/datetime";

  const { weather } = $props<{ weather: WeatherCard }>();

  const formattedTime = $derived(formatTime(weather.timestamp));
</script>

<div class="grid grid-cols-3 items-center relative">
  <!-- Left section - City and time -->
  <div class="flex flex-col">
    <h2 class="text-3xl font-semibold text-white">{weather.city}</h2>
    <p class="text-blue-300 text-sm mt-1">{formattedTime}</p>
  </div>

  <!-- Middle section - Weather icon -->
  <div class="flex justify-center">
    <WeatherIcon code={weather.code} />
  </div>

  <!-- Right section - Temperature and conditions -->
  <div class="flex flex-col items-end">
    <div class="text-3xl font-bold text-white">{weather.temperature}°C</div>
    <p class="text-blue-300 text-sm">
      {weather.humidity}%, {weather.precipitation}mm
    </p>
  </div>
</div>
