<script lang="ts">
  import WeatherIcon from "$comp/weather/WeatherIcon.svelte";

  interface WeatherProps {
    cityName: string;
    timestamp: string;
    temperature: number;
    humidity: number;
    precipitation: number;
    weatherCode: number;
  }

  let {
    cityName,
    timestamp,
    temperature,
    humidity,
    precipitation,
    weatherCode,
  }: WeatherProps = $props();

  const date = new Date(timestamp);
  const formattedTime = $derived(
    `${date.toLocaleDateString("en-US", { weekday: "long" })}, ${date.toLocaleTimeString(
      "en-US",
      {
        hour: "2-digit",
        minute: "2-digit",
        hour12: false,
      },
    )}`,
  );
</script>

<div class="flex justify-between items-start relative">
  <div>
    <h2 class="text-2xl font-semibold text-blue-100">{cityName}</h2>
    <p class="text-blue-300 text-sm mt-1">{formattedTime}</p>
  </div>
  <div class="absolute left-1/2 -translate-x-1/2 top-1/2 -translate-y-1/2">
    <WeatherIcon code={weatherCode} />
  </div>
  <div class="text-right">
    <div class="text-4xl font-bold text-blue-100">{temperature}°C</div>
    <p class="text-blue-300 text-sm">{humidity}%, {precipitation}mm</p>
  </div>
</div>