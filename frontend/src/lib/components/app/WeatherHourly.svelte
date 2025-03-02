<script lang="ts">
  import type { HourlyForecast } from "$type/Forecast";

  type ProcessedForecast = {
    time: string;
    temperature: number;
    humidity: number;
    weatherCode: number;
    precipitation: number;
  };

  let { hourlyForecast }: {hourlyForecast: HourlyForecast} = $props<{ hourlyForecast: HourlyForecast }>();

  function formatTime(isoString: string): string {
    const date = new Date(isoString);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false });
  }

  function formatTemp(temp: number): string {
    return `${Math.round(temp)}°C`;
  }

  function formatHumidity(humidity: number): string {
    return `${humidity}%`;
  }

  function formatPrecipitation(precipitation: number): string {
    return `${precipitation}mm`;
  }


  function getWeatherIcon(code: number): string {
    if (code === 0) return '☀️'; // Clear sky
    if (code === 1) return '🌤️'; // Mainly clear
    if (code === 2) return '⛅'; // Partly cloudy
    if (code === 3) return '☁️'; // Overcast
    if (code >= 45 && code <= 48) return '🌫️'; // Fog
    if (code >= 51 && code <= 55) return '🌧️'; // Drizzle
    if (code >= 61 && code <= 65) return '🌧️'; // Rain
    if (code >= 71 && code <= 75) return '❄️'; // Snow
    if (code >= 80 && code <= 82) return '🌦️'; // Showers
    if (code >= 95) return '⛈️'; // Thunderstorm
    return '🌡️'; // Default
  }

  const forecastData = $derived(processWeatherData(hourlyForecast));

  function processWeatherData(data?: HourlyForecast): ProcessedForecast[] {
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

<style>
  .weather-table {
    border-collapse: collapse;
    color: var(--color-table-text);
  }

  .weather-table tr:nth-child(odd) td {
    background-color: var(--color-table-row-odd);
  }

  .weather-table tr:nth-child(even) td {
    background-color: var(--color-table-row-even);
  }

  .weather-table td {
    padding: 0.75rem;
  }

  .weather-table .temperature {
    font-weight: 600;
    color: var(--color-table-accent);
  }

  .weather-icon {
    font-size: 1.5rem;
  }
</style>

<div>
  {#if forecastData.length > 0}
    <table class="w-full weather-table">
      <tbody>
        {#each forecastData as forecast}
          <tr>
            <td>{formatTime(forecast.time)}</td>
            <td>
              <div class="weather-icon">{getWeatherIcon(forecast.weatherCode)}</div>
            </td>
            <td class="temperature">{formatTemp(forecast.temperature)}</td>
            <td>{formatHumidity(forecast.humidity)}</td>
            <td>{formatPrecipitation(forecast.precipitation)}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>Loading weather data...</p>
  {/if}
</div>