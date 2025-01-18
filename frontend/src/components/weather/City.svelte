<script lang="ts">
  interface WeatherProps {
    cityName: string;
    temperature: number;
    condition: string;
    timestamp: string;
  }

  let { cityName, temperature, condition, timestamp }: WeatherProps = $props();
</script>

<div class="flex justify-between items-start">
  <div>
    <h2 class="text-2xl font-semibold text-blue-100">{cityName}</h2>
    <p class="text-blue-300 text-sm mt-1">{timestamp}</p>
  </div>
  <div class="text-right">
    <div class="text-4xl font-bold text-blue-100">{temperature}°C</div>
    <p class="text-blue-300 text-sm">{condition}</p>
  </div>
</div>