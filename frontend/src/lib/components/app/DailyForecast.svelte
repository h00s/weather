<script lang="ts">
  import type { DailyForecast } from "$type/Forecast";
  import { formatTemperature, formatPrecipitation, weatherIcon } from "$util/weather";
  import { formatDateWithWeekday } from "$util/datetime";

  type ProcessedDailyForecast = {
    date: string;
    maxTemperature: number;
    minTemperature: number;
    weatherCode: number;
    precipitation: number;
  };

  const { data }: { data: DailyForecast } = $props<{ data: DailyForecast }>();

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
                {weatherIcon(forecast.weatherCode)}
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
