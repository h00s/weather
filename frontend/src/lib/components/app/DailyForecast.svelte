<script lang="ts">
  import type { DailyForecast } from "$type/Forecast";
  import { formatTemperature, formatPrecipitation } from "$util/weather";
  import { formatDateWithWeekday } from "$util/datetime";

  type ProcessedDailyForecast = {
    date: string;
    maxTemperature: number;
    minTemperature: number;
    weatherCode: number;
    precipitation: number;
  };

  const { data }: { data: DailyForecast } = $props<{ data: DailyForecast }>();

  function getWeatherIcon(code: number): string {
    if (code === 0) return "☀️"; // Clear sky
    if (code === 1) return "🌤️"; // Mainly clear
    if (code === 2) return "⛅"; // Partly cloudy
    if (code === 3) return "☁️"; // Overcast
    if (code >= 45 && code <= 48) return "🌫️"; // Fog
    if (code >= 51 && code <= 55) return "🌧️"; // Drizzle
    if (code >= 61 && code <= 65) return "🌧️"; // Rain
    if (code >= 71 && code <= 75) return "❄️"; // Snow
    if (code >= 80 && code <= 82) return "🌦️"; // Showers
    if (code >= 95) return "⛈️"; // Thunderstorm
    return "🌡️"; // Default
  }

  const forecastData = $derived(processDailyForecast(data));

  function processDailyForecast(
    data?: DailyForecast,
  ): ProcessedDailyForecast[] {
    if (!data || !data.time || !data.data) return [];

    const processed: ProcessedDailyForecast[] = [];

    const displayDays = Math.min(data.time.length, 7);

    for (let i = 0; i < displayDays; i++) {
      processed.push({
        date: data.time[i],
        maxTemperature: data.data.temperature_2m_max[i],
        minTemperature: data.data.temperature_2m_min[i],
        weatherCode: data.data.weather_code[i],
        precipitation: data.data.precipitation_sum[i],
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
            <td>{formatDateWithWeekday(forecast.date)}</td>
            <td>
              <div class="weather-icon">
                {getWeatherIcon(forecast.weatherCode)}
              </div>
            </td>
            <td class="primary">
              {formatTemperature(forecast.maxTemperature)} / {formatTemperature(
                forecast.minTemperature,
              )}
            </td>
            <td class="secondary"
              >{formatPrecipitation(forecast.precipitation)}</td
            >
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>Loading forecast data...</p>
  {/if}
</div>
