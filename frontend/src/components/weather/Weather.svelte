<script lang="ts">
  interface WeatherProps {
    cityName: string;
    timestamp: string;
    temperature: number;
    humidity: number;
    precipitation: number;
  }

  let {
    cityName,
    timestamp,
    temperature,
    humidity,
    precipitation,
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

<div class="flex justify-between items-start">
  <div>
    <h2 class="text-2xl font-semibold text-blue-100">{cityName}</h2>
    <p class="text-blue-300 text-sm mt-1">{formattedTime}</p>
  </div>
  <div class="text-right">
    <div class="text-4xl font-bold text-blue-100">{temperature}°C</div>
    <p class="text-blue-300 text-sm">{humidity}%, {precipitation}mm</p>
  </div>
</div>
