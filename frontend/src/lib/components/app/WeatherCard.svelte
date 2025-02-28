<script lang="ts">
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
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="48"
      height="48"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      class="lucide-icon lucide lucide-cloudy text-blue-100"
    >
      <path d="M17.5 21H9a7 7 0 1 1 6.71-9h1.79a4.5 4.5 0 1 1 0 9Z"></path>
      <path d="M22 10a3 3 0 0 0-3-3h-2.207a5.502 5.502 0 0 0-10.702.5"></path>
    </svg>
  </div>

  <!-- Right section - Temperature and conditions -->
  <div class="flex flex-col items-end">
    <div class="text-3xl font-bold text-white">{weather.temperature}°C</div>
    <p class="text-blue-300 text-sm">
      {weather.humidity}%, {weather.precipitation}mm
    </p>
  </div>
</div>
