<script lang="ts">
  import type { HourlyForecast } from "$type/Forecast";
  import { formatTime } from "$util/datetime";
  import { formatTemperature, formatHumidity, formatPrecipitation, weatherIcon } from "$util/weather";

  type ProcessedForecast = {
    time: string;
    temperature: number;
    humidity: number;
    weatherCode: number;
    precipitation: number;
  };

  const { data }: { 
    data: HourlyForecast;
  } = $props<{ 
    data: HourlyForecast;
  }>();

  const forecastData = $derived(processHourlyForecast(data));

  function processHourlyForecast(data?: HourlyForecast): ProcessedForecast[] {
    if (!data || !data.time || !data.data) return [];
    
    const processed: ProcessedForecast[] = [];
    
    const now = new Date();
    const currentHour = now.getHours();
    
    const startIndex = data.time.findIndex(timeStr => {
      const timeDate = new Date(timeStr);
      return timeDate.getHours() >= currentHour;
    });
    
    const effectiveStartIndex = startIndex >= 0 ? startIndex : 0;
    
    for (let i = effectiveStartIndex; i < effectiveStartIndex + 12 && i < data.time.length; i++) {
      processed.push({
        time: data.time[i],
        temperature: data.data.temperature_2m[i],
        humidity: data.data.relative_humidity_2m[i],
        weatherCode: data.data.weather_code[i],
        precipitation: data.data.precipitation[i]
      });
    }
    
    return processed;
  }
</script>

<div>
  {#if forecastData.length > 0}
    <table class="weather-table">
      <tbody>
        {#each forecastData as forecast}
          <tr>
            <td>{formatTime(forecast.time)}</td>
            <td>
              <div class="weather-icon">{weatherIcon(forecast.weatherCode)}</div>
            </td>
            <td class="primary">{formatTemperature(forecast.temperature)}</td>
            <td class="secondary">{formatHumidity(forecast.humidity)}</td>
            <td class="secondary">{formatPrecipitation(forecast.precipitation)}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>Loading weather data...</p>
  {/if}
</div>